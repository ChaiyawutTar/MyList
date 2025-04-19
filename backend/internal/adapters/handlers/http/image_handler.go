// internal/adapters/handlers/http/image_handler.go
package http

import (
    "fmt"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/ChaiyawutTar/MyList/internal/core/ports"

	"strings"
)

type ImageHandler struct {
    imageRepository ports.ImageRepository
}

func NewImageHandler(imageRepository ports.ImageRepository) *ImageHandler {
    return &ImageHandler{
        imageRepository: imageRepository,
    }
}


// internal/adapters/handlers/http/image_handler.go
func (h *ImageHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
    // Get image ID from URL
    imageID := chi.URLParam(r, "id")
    if imageID == "" {
        http.Error(w, "Image ID is required", http.StatusBadRequest)
        return
    }

    // Log the request
    fmt.Printf("Serving image with ID: %s\n", imageID)

    // Generate ETag based on image ID
    etag := fmt.Sprintf("\"img-%s\"", imageID)
    
    // Check if the client has the image cached
    if match := r.Header.Get("If-None-Match"); match != "" {
        if strings.Contains(match, etag) {
            w.WriteHeader(http.StatusNotModified)
            return
        }
    }

    // Get image data from repository with error handling
    imageData, contentType, err := h.imageRepository.Get(r.Context(), imageID)
    if err != nil {
        if strings.Contains(err.Error(), "not found") {
            fmt.Printf("Image not found: %s, error: %v\n", imageID, err)
            http.Error(w, "Image not found", http.StatusNotFound)
        } else {
            fmt.Printf("Error retrieving image %s: %v\n", imageID, err)
            http.Error(w, "Failed to retrieve image", http.StatusInternalServerError)
        }
        return
    }

    // Check if we actually got data
    if len(imageData) == 0 {
        fmt.Printf("Image data is empty for ID: %s\n", imageID)
        http.Error(w, "Image data is empty", http.StatusNotFound)
        return
    }

    // Set strong caching headers
    w.Header().Set("ETag", etag)
    w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache for 1 year
    w.Header().Set("Content-Type", contentType)
    w.Header().Set("Content-Length", fmt.Sprintf("%d", len(imageData)))
    
    // Write image data to response
    bytesWritten, err := w.Write(imageData)
    if err != nil {
        fmt.Printf("Error writing image data to response: %v\n", err)
        return
    }
    
    fmt.Printf("Successfully served image %s (%d bytes)\n", imageID, bytesWritten)
}