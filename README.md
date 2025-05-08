# FileHttp-Go

FileHttp-Go is a lightweight HTTP server written in Go for securely uploading and serving JSON files. It includes built-in support for file validation, concurrency management, and temporary file handling.

## Features

- **File Upload**: Users can upload `.json` files via HTTP POST requests.
- **File Validation**: Ensures uploaded files are valid `.json` files and verifies filenames.
- **Concurrency**: Limits the number of concurrent uploads using a semaphore.
- **Temporary File Handling**: Uploaded files are initially saved as temporary files and then moved to the target directory.
- **File Serving**: Serves uploaded files through an HTTP file server for easy access.

## Getting Started

### Prerequisites

- Go 1.16 or later must be installed.
- Basic knowledge of HTTP and Go programming.

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/the0cp/fileHttp-go.git
   cd fileHttp-go
   ```

2. Build the project:
   ```bash
   go build
   ```

### Usage

Run the server with the following command:
```bash
./fileHttp-go -dir <upload-directory> -port <port-number>
```
- Replace `<upload-directory>` with the directory where you want to save uploaded files (default: `.`).
- Replace `<port-number>` with the desired port number (default: `8080`).

Example:
```bash
./fileHttp-go -dir uploads -port 8080
```

### Endpoints

1. **File Upload**
   - **URL**: `/upload`
   - **Method**: `POST`
   - **Query Parameter**: `filename` (required, must be a `.json` file)
   - **Body**: Raw JSON file
   - **Example**:
     ```bash
     curl -X POST -F "filename=test.json" --data-binary @test.json http://localhost:8080/upload
     ```

2. **File Serving**
   - Access uploaded files through the root endpoint /.
   - Example: http://localhost:8080/<filename>.json

### Configuration Options

- **Upload Directory**: Specify where uploaded files are stored using the `-dir` flag.
- **Port**: Set the port for the HTTP server using the `-port` flag.

### Limitations

- Only `.json` files are accepted for upload.
- File size is limited to 5 MB per upload. Modify the `maxUploadSize` constant in `server.go` to change this limit.

## Development

### Code Structure

- **`server.go`**: Main application logic, including file upload handling and HTTP server configuration.

### Running Locally

To run the server locally for development purposes:
```bash
go run server.go -dir ./uploads -port 8080
```

## License

This project is licensed under the GNU General Public License v3. See the [LICENSE](LICENSE) file for details.
