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

#### Example:

```bash
./fileHttp-go -dir uploads -port 8443 -ca ca.crt
```

## License

This project is licensed under the GNU General Public License v3. See the [LICENSE](LICENSE) file for details.
