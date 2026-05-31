package services

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/rpdg/vercel_blob"
)

func UploadFileToVercelBlob(file *multipart.FileHeader, prefix string) (string, error) {
	client := vercel_blob.NewVercelBlobClient()

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %w", err)
	}

	filename := fmt.Sprintf("%s_%d_%s", prefix, time.Now().Unix(), filepath.Base(file.Filename))

	log.Printf("Uploading file '%s' to Vercel Blob...", filename)

	putResult, err := client.Put(filename, bytes.NewReader(fileBytes), vercel_blob.PutCommandOptions{
		ContentType: file.Header.Get("Content-Type"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload blob '%s': %w", filename, err)
	}

	log.Printf("File '%s' uploaded to Vercel Blob. URL: %s", filename, putResult.URL)

	return putResult.URL, nil
}