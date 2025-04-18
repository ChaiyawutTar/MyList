package file

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type imageRepository struct {
	uploadDir string
}

func NewImageRepository(uploadDir string) *imageRepository {
	// Create upload directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		panic(err)
	}

	return &imageRepository{
		uploadDir: uploadDir,
	}
}

func (r *imageRepository) Save(ctx context.Context, file multipart.File, filename string) (string, error) {
	// Create destination file
	filePath := filepath.Join(r.uploadDir, filename)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy file content
	if _, err := io.Copy(dst, file); err != nil {
		return "", err
	}

	// Return relative path
	return filepath.Join("uploads", filename), nil
}

func (r *imageRepository) Delete(ctx context.Context, path string) error {
	// Get absolute path
	absPath := filepath.Join(".", path)
	return os.Remove(absPath)
}