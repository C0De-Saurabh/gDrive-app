package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Unable to listen on port 8080: %v", err)
	}
	defer listener.Close()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the authorization code: \n%v\n", authURL)

	codeChan := make(chan string)
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			if code != "" {
				codeChan <- code
				fmt.Fprintf(w, "Authorization code received. You can now close this window.")
			} else {
				fmt.Fprintf(w, "No authorization code received.")
			}
		})
		http.Serve(listener, nil)
	}()

	authCode := <-codeChan

	tok, err := config.Exchange(context.Background(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieve a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Save a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to create token file: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

// ListAllFiles retrieves all files in the user's Drive.
func ListAllFiles(srv *drive.Service) ([]*drive.File, error) {
	var files []*drive.File
	pageToken := ""
	for {
		r, err := srv.Files.List().PageSize(1000).Fields("nextPageToken, files(id, name, mimeType, size, md5Checksum, parents)").PageToken(pageToken).Do()
		if err != nil {
			return nil, err
		}
		files = append(files, r.Files...)
		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return files, nil
}

// GetFilePaths retrieves the full paths of the files from their parent IDs.
func GetFilePaths(srv *drive.Service, files []*drive.File) map[string]string {
	paths := make(map[string]string)
	for _, file := range files {
		if len(file.Parents) > 0 {
			path := getFullPath(srv, file.Parents[0])
			paths[file.Id] = path
		} else {
			paths[file.Id] = "root"
		}
	}
	return paths
}

// getFullPath recursively retrieves the full path of a folder by its ID.
func getFullPath(srv *drive.Service, folderID string) string {
	folder, err := srv.Files.Get(folderID).Fields("id, name, parents").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve folder: %v", err)
	}
	if len(folder.Parents) > 0 {
		return getFullPath(srv, folder.Parents[0]) + "/" + folder.Name
	}
	return folder.Name
}

// FindDuplicates identifies duplicate files by md5Checksum.
func FindDuplicates(files []*drive.File) map[string][]*drive.File {
	duplicates := make(map[string][]*drive.File)
	hashMap := make(map[string]*drive.File)

	for _, file := range files {
		if existing, found := hashMap[file.Md5Checksum]; found {
			duplicates[file.Md5Checksum] = append(duplicates[file.Md5Checksum], file)
			if len(duplicates[file.Md5Checksum]) == 1 {
				duplicates[file.Md5Checksum] = append(duplicates[file.Md5Checksum], existing)
			}
		} else {
			hashMap[file.Md5Checksum] = file
		}
	}

	return duplicates
}

// PrintDuplicates prints the duplicate files along with their paths.
func PrintDuplicates(srv *drive.Service, duplicates map[string][]*drive.File, paths map[string]string) {
	fmt.Println("Duplicate files:")
	if len(duplicates) == 0 {
		fmt.Println("No duplicate files found.")
	} else {
		for hash, files := range duplicates {
			fmt.Printf("Hash: %s\n", hash)
			for _, file := range files {
				fmt.Printf("  ID: %s, Name: %s, Size: %d, Path: %s\n", file.Id, file.Name, file.Size, paths[file.Id])
			}
		}
	}
}

func main() {
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	files, err := ListAllFiles(srv)
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}

	paths := GetFilePaths(srv, files)
	duplicates := FindDuplicates(files)
	PrintDuplicates(srv, duplicates, paths)
}
