// internal/adapters/repository/postgres/image_repository.go
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"time"
	"net/http"
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
    // Read file content
    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        return "", err
    }

    // Detect content type
    contentType := http.DetectContentType(fileBytes)

    // Store in database
    var id int
    err = r.db.QueryRowContext(
        ctx,
        "INSERT INTO images (filename, data, content_type, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
        filename, fileBytes, contentType, time.Now(),
    ).Scan(&id)
    if err != nil {
        return "", err
    }

    // Return the ID as a string reference
    return fmt.Sprintf("%d", id), nil
}

func (r *imageRepository) Delete(ctx context.Context, imageID string) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM images WHERE id = $1", imageID)
	return err
}

func (r *imageRepository) Get(ctx context.Context, imageID string) ([]byte, string, error) {
    fmt.Printf("Getting image with ID: %s\n", imageID)
    
    var data []byte
    var contentType sql.NullString
    
    err := r.db.QueryRowContext(
        ctx, 
        "SELECT data, content_type FROM images WHERE id = $1", 
        imageID,
    ).Scan(&data, &contentType)
    
    if err != nil {
        fmt.Printf("Error retrieving image from database: %v\n", err)
        return nil, "", err
    }
    
    // Use a default content type if it's NULL in the database
    contentTypeStr := "image/jpeg"
    if contentType.Valid {
        contentTypeStr = contentType.String
    }
    
    return data, contentTypeStr, nil
}