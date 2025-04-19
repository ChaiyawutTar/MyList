// internal/adapters/repositories/postgres/image_repository.go
package postgres  // Make sure this is postgres, not file

import (
    "context"
    "database/sql"
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "time"
    "bytes"
	"strconv"
    "github.com/ChaiyawutTar/MyList/internal/core/ports"
)

type imageRepository struct {
    db *sql.DB
}

func NewImageRepository(db *sql.DB) ports.ImageRepository {
    // Create the images table if it doesn't exist
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS images (
            id SERIAL PRIMARY KEY,
            filename TEXT NOT NULL,
            data BYTEA NOT NULL,
            content_type TEXT,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
    if err != nil {
        panic(err)
    }

    return &imageRepository{
        db: db,
    }
}

func (r *imageRepository) Save(ctx context.Context, file multipart.File, filename string) (string, error) {
    // Read the first few bytes to detect content type
    buffer := make([]byte, 512)
    _, err := file.Read(buffer)
    if err != nil && err != io.EOF {
        return "", fmt.Errorf("failed to read file header: %w", err)
    }
    
    // Detect content type
    contentType := http.DetectContentType(buffer)
    
    // Reset file pointer to beginning
    _, err = file.Seek(0, io.SeekStart)
    if err != nil {
        return "", fmt.Errorf("failed to reset file pointer: %w", err)
    }
    
    // Read the entire file into memory
    var fileBuffer bytes.Buffer
    _, err = io.Copy(&fileBuffer, file)
    if err != nil {
        return "", fmt.Errorf("failed to read file: %w", err)
    }
    
    fileBytes := fileBuffer.Bytes()
    
    // Begin a transaction
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return "", fmt.Errorf("failed to begin transaction: %w", err)
    }
    
    // Prepare to rollback in case of error
    defer func() {
        if err != nil {
            tx.Rollback()
        }
    }()
    
    // Store in database
    var id int
    err = tx.QueryRowContext(
        ctx,
        "INSERT INTO images (filename, data, content_type, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
        filename, fileBytes, contentType, time.Now(),
    ).Scan(&id)
    if err != nil {
        return "", fmt.Errorf("failed to insert image: %w", err)
    }
    
    // Commit the transaction
    err = tx.Commit()
    if err != nil {
        return "", fmt.Errorf("failed to commit transaction: %w", err)
    }
    
    fmt.Printf("Successfully saved image with ID: %d, size: %d bytes\n", id, len(fileBytes))
    
    // Return the ID as a string reference
    return fmt.Sprintf("%d", id), nil
}

func (r *imageRepository) Delete(ctx context.Context, imageID string) error {
    fmt.Printf("Deleting image with ID: %s\n", imageID)
    
    result, err := r.db.ExecContext(ctx, "DELETE FROM images WHERE id = $1", imageID)
    if err != nil {
        return fmt.Errorf("failed to delete image: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("image with ID %s not found", imageID)
    }
    
    fmt.Printf("Successfully deleted image with ID: %s\n", imageID)
    return nil
}

func (r *imageRepository) Get(ctx context.Context, imageID string) ([]byte, string, error) {
    // Convert string ID to int
    id, err := strconv.Atoi(imageID)
    if err != nil {
        return nil, "", fmt.Errorf("invalid image ID: %w", err)
    }

    // Query the database
    var data []byte
    var contentType string
    var filename string
    
    query := "SELECT data, content_type, filename FROM images WHERE id = $1"
    err = r.db.QueryRowContext(ctx, query, id).Scan(&data, &contentType, &filename)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, "", fmt.Errorf("image not found: %s", imageID)
        }
        return nil, "", fmt.Errorf("error querying image: %w", err)
    }

    // Validate the data
    if len(data) == 0 {
        return nil, "", fmt.Errorf("image data is empty for ID: %s", imageID)
    }

    // If content type is empty, try to detect it
    if contentType == "" {
        contentType = http.DetectContentType(data[:512])
    }

    return data, contentType, nil
}