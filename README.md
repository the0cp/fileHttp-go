# FileHttp-Go

FileHttp-Go is a lightweight HTTP server written in Go for securely uploading and serving JSON files. It includes built-in support for file validation, concurrency management, and temporary file handling.

## Features

- **File Upload**: Users can upload `.json` files via HTTP POST requests.
- **File Validation**: Ensures uploaded files are valid `.json` files and verifies filenames.
- **Concurrency**: Limits the number of concurrent uploads using a semaphore.
- **Temporary File Handling**: Uploaded files are initially saved as temporary files and then moved to the target directory.
- **File Serving**: Serves uploaded files through an HTTP file server for easy access.
- **Mutual TLS (mTLS) Authentication**: The `/upload` endpoint requires clients to provide a trusted certificate, enforcing device whitelisting.

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

Run the server with the following command (specifying the CA certificate for mTLS):

```bash
./fileHttp-go -dir <upload-directory> -port <port-number> -ca <path-to-ca.crt>
```

- `<upload-directory>`: Directory to save uploaded files (default: `.`)
- `<port-number>`: Port to listen on (default: `8080`)
- `<path-to-ca.crt>`: Path to the CA certificate used for client certificate verification (required for mTLS)

#### Example:

```bash
./fileHttp-go -dir uploads -port 8443 -ca ca.crt
```

### Endpoints

1. **File Upload (Requires client certificate)**
    - **URL**: `/upload`
    - **Method**: `POST`
    - **Query Parameter**: `filename` (required, must be a `.json` file)
    - **Body**: Raw JSON file
    - **Example (with client certs)**:
        ```bash
        curl --cert client.crt --key client.key --cacert ca.crt -X POST -F "filename=test.json" --data-binary @test.json https://localhost:8443/upload
        ```

2. **File Serving**
    - Access uploaded files through the root endpoint `/`.
    - Example: https://localhost:8443/<filename>.json

### Configuration Options

- **Upload Directory**: Specify where uploaded files are stored using the `-dir` flag.
- **Port**: Set the port for the HTTP server using the `-port` flag.
- **CA Certificate**: Set the CA certificate for mTLS with the `-ca` flag.

## License

This project is licensed under the GNU General Public License v3. See the [LICENSE](LICENSE) file for details.
