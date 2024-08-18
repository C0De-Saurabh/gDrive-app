# Google Drive CLI Application

## Overview

The Google Drive CLI Application is a command-line tool developed in Golang that allows users to interact with Google Drive. The application enables functionalities such as listing, navigating, and identifying duplicate files based on their content hashes.

## Features

- **OAuth 2.0 Authentication:** Securely access Google Drive APIs.
- **File Listing and Navigation:** View all files and directories within Google Drive, with hierarchical structure display.
- **Duplicate File Detection:** Identify and list duplicate files using Google Drive's `md5Checksum`.
- **Enhanced User Experience:** Display paths and metadata of duplicate files for efficient management.
- **Security:** Manage and ignore sensitive files and build artifacts.

# Installation

## Prerequisites

Before you begin, ensure you have the following:

- **Go (Golang) 1.16 or later**: You can download and install Go from [golang.org](https://golang.org/dl/).
- **Google Cloud Platform Project**: Create a project, enable the Google Drive API, and download the `credentials.json` file from the Google Cloud Console.
- **Git**: Make sure Git is installed on your system.

## Steps

1. **Clone the Repository**

   First, clone the repository to your local machine:
   ```bash
   git clone https://github.com/yourusername/google-drive-cli.git
   cd google-drive-cli
   ```
2. **Install Dependencies**

   Download and install the required Go modules:
   ```bash
   go mod tidy
   ```
3. **Configure OAuth 2.0**

   - Download the `credentials.json` file from your Google Cloud Console, where you have enabled the Google Drive API for your project.
   - Place the `credentials.json` file in the root directory of the project. This file is required for OAuth 2.0 authentication and will allow the application to securely access your Google Drive.
4. **Build the Application**

   Compile the application using Go to create the executable:
   ```bash
   go build -o gdrivecli
    ```
5. **Run the Application**

   Start the application and initiate the OAuth 2.0 authentication process:
   ```bash
   ./gdrivecli auth
   ```
