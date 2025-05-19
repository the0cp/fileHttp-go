package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const maxUploadSize = 5 << 20 // 5 MB

var uploadDir string

var semaphore = make(chan struct{}, 1000)

func main() {
	dir := flag.String("dir", ".", "Directory to save uploaded files")
	port := flag.String("port", "8080", "Port to listen on")
	caPath := flag.String("ca", "ca.crt", "CA Certificate Path")
	flag.Parse()

	uploadDir = *dir

	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		log.Printf("Directory %s does not exist, creating it...", uploadDir)
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			log.Fatalf("Failed to create upload directory: %v", err)
		}
		log.Printf("Directory %s created successfully", uploadDir)
	} else if err != nil {
		log.Fatalf("Failed to check upload directory: %v", err)
	}

	fileServer := http.FileServer(http.Dir(uploadDir))
	http.Handle("/", fileServer)
	http.Handle("/upload", mTLSAuthMidware(uploadHandler))

	tlsConfig, err := loadTLSConfig(*caPath)
	if err != nil {
		log.Fatalf("TLS config error: %v", err)
	}

	server := &http.Server{
		Addr:      ":" + *port,
		TLSConfig: tlsConfig,
	}

	log.Printf("Server starting on port %s...", *port)
	log.Printf("Files will be saved in %s", uploadDir)

	if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	semaphore <- struct{}{}
	defer func() {
		<-semaphore
		log.Printf("Semaphore released, current capacity: %d", cap(semaphore)-len(semaphore))
	}()

	if r.ContentLength > maxUploadSize {
		http.Error(w, fmt.Sprintf("File size exceeds limit of %d bytes", maxUploadSize), http.StatusBadRequest)
		log.Printf("Upload request exceeds maximum size of %d bytes", maxUploadSize)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	filename := r.URL.Query().Get("filename")
	if filename == "" {
		http.Error(w, "Filename is required in query parameter", http.StatusBadRequest)
		log.Println("Upload request received without filename query parameter")
		return
	}

	baseFilename := filepath.Base(filename)
	if baseFilename == "" || baseFilename == "." || baseFilename == ".." || strings.Contains(baseFilename, "..") || filepath.IsAbs(baseFilename) {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		log.Printf("Received invalid filename: %s", filename)
		return
	}

	fileExt := strings.ToLower(filepath.Ext(baseFilename))
	if fileExt != ".json" {
		http.Error(w, fmt.Sprintf("Invalid file extension '%s'. Only .json files are allowed.", fileExt), http.StatusBadRequest)
		log.Printf("Received file with invalid extension: %s (original filename: %s)", fileExt, filename)
		return
	}

	tempFile, err := os.CreateTemp("", "upload-")
	if err != nil {
		http.Error(w, "Server error: Failed to create temporary file", http.StatusInternalServerError)
		log.Printf("Failed to create temporary file: %v", err)
		return
	}
	tempFilePath := tempFile.Name()
	defer cleanupTempFiles(tempFile, tempFilePath)

	log.Printf("Starting copy to temporary file: %s", tempFilePath)
	n, err := io.Copy(tempFile, r.Body)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFilePath)
		http.Error(w, fmt.Sprintf("Failed to read or save temporary file: %v", err), http.StatusInternalServerError)
		log.Printf("Failed to copy request body to temporary file %s: %v", tempFilePath, err)
		return
	}
	log.Printf("Copied %d bytes to temporary file: %s", n, tempFilePath)

	go saveFile(tempFilePath, baseFilename, uploadDir)

	fmt.Fprintf(w, "File '%s' received and is being saved in background.", baseFilename)
	log.Printf("File '%s' received, background save initiated (temp file: %s).", baseFilename, tempFilePath)
}

func saveFile(tempFilePath string, finalFilename string, destDir string) {
	defer func() {
		err := os.Remove(tempFilePath)
		if err != nil {
			log.Printf("Background cleanup error: Failed to remove temporary file %s: %v", tempFilePath, err)
		} else {
			log.Printf("Background cleanup successful: Removed temporary file %s", tempFilePath)
		}
	}()
	log.Printf("Background save started: From temporary file %s to %s/%s", tempFilePath, destDir, finalFilename)

	tempFile, err := os.Open(tempFilePath)
	if err != nil {
		log.Printf("Background save error: Failed to open temporary file %s: %v", tempFilePath, err)
		return
	}
	defer tempFile.Close()

	dstPath := filepath.Join(destDir, finalFilename)

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Printf("Background save error: Failed to create destination file '%s': %v", dstPath, err)
		return
	}
	defer dst.Close()

	log.Printf("Background save: Copying content from %s to %s", tempFilePath, dstPath)
	if _, err := io.Copy(dst, tempFile); err != nil {
		log.Printf("Background save error: Failed to copy file content to '%s': %v", dstPath, err)
		return
	}

	log.Printf("Background save successful: File '%s' saved to '%s'", finalFilename, dstPath)
}

func cleanupTempFiles(tempFile *os.File, tempFilePath string) {
	tempFile.Close()
	err := os.Remove(tempFilePath)
	if err != nil {
		log.Printf("Failed to remove temporary file %s: %v", tempFilePath, err)
	} else {
		log.Printf("Temporary file %s removed successfully", tempFilePath)
	}
}
