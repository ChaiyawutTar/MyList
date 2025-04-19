// internal/adapters/handlers/http/image_handler.go
package http

import (
    "fmt"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/ChaiyawutTar/MyList/internal/core/ports"
)

type ImageHandler struct {
    imageRepository ports.ImageRepository
}

func NewImageHandler(imageRepository ports.ImageRepository) *ImageHandler {
    return &ImageHandler{
        imageRepository: imageRepository,
    }
}


func (h *ImageHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
    // Get image ID from URL
    imageID := chi.URLParam(r, "id")
    if imageID == "" {
        http.Error(w, "Image ID is required", http.StatusBadRequest)
        return
    }

    fmt.Printf("Serving image with ID: %s\n", imageID)

    // Get image data from repository
    imageData, contentType, err := h.imageRepository.Get(r.Context(), imageID)
    if err != nil {
        fmt.Printf("Error retrieving image %s: %v\n", imageID, err)
        http.Error(w, fmt.Sprintf("Failed to retrieve image: %v", err), http.StatusInternalServerError)
        return
    }

    fmt.Printf("Successfully retrieved image %s with content type %s and size %d bytes\n", 
        imageID, contentType, len(imageData))

    // Set content type
    w.Header().Set("Content-Type", contentType)
    w.Header().Set("Content-Length", fmt.Sprintf("%d", len(imageData)))
    w.Header().Set("Cache-Control", "public, max-age=86400") // Cache for 24 hours
    
    // Write image data to response
    _, err = w.Write(imageData)
    if err != nil {
        fmt.Printf("Error writing image data to response: %v\n", err)
    }
}