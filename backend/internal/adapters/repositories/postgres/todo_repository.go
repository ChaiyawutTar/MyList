package postgres

import (
	"context"
	"database/sql"

	"errors"
	"time"
	"fmt"

	"github.com/ChaiyawutTar/MyList/internal/core/domain"
)

type todoRepository struct {
	db *sql.DB
}

func (r *todoRepository) verifyDatabaseSetup(ctx context.Context) error {
    // 1. Check database connection
    if err := r.db.PingContext(ctx); err != nil {
        return fmt.Errorf("database connection error: %w", err)
    }

    // 2. Verify table exists and has correct structure
    verifyQuery := `
        SELECT column_name 
        FROM information_schema.columns 
        WHERE table_name = 'todos' 
        ORDER BY ordinal_position;`

    rows, err := r.db.QueryContext(ctx, verifyQuery)
    if err != nil {
        return fmt.Errorf("error verifying table structure: %w", err)
    }
    defer rows.Close()

    columns := make([]string, 0)
    for rows.Next() {
        var columnName string
        if err := rows.Scan(&columnName); err != nil {
            return fmt.Errorf("error scanning column name: %w", err)
        }
        columns = append(columns, columnName)
    }

    fmt.Printf("Found columns in todos table: %v\n", columns)

    // 3. Verify we can execute a simple query
    testQuery := `SELECT COUNT(*) FROM todos WHERE user_id = $1;`
    var count int
    if err := r.db.QueryRowContext(ctx, testQuery, 1).Scan(&count); err != nil {
        return fmt.Errorf("error executing test query: %w", err)
    }

    return nil
}

func NewTodoRepository(db *sql.DB) *todoRepository {
    repo := &todoRepository{db: db}
    
    // Verify database setup
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := repo.verifyDatabaseSetup(ctx); err != nil {
        panic(fmt.Sprintf("database setup verification failed: %v", err))
    }
    
    return repo
}


func (r *todoRepository) FindAll(ctx context.Context, userID int) ([]domain.Todo, error) {
    // 1. First, let's verify the table structure
    verifyQuery := `SELECT COUNT(*) FROM information_schema.columns WHERE table_name = 'todos';`
    var columnCount int
    err := r.db.QueryRowContext(ctx, verifyQuery).Scan(&columnCount)
    if err != nil {
        return nil, fmt.Errorf("error verifying table structure: %w", err)
    }
    fmt.Printf("Number of columns in todos table: %d\n", columnCount)

    // 2. Use a simpler query first to debug
    query := `
        SELECT id, user_id, title, description, status, COALESCE(image_id, ''), created_at, updated_at 
        FROM todos 
        WHERE user_id = $1 
        ORDER BY created_at DESC`

    // 3. Print the query and userID for debugging
    fmt.Printf("Executing query: %s with userID: %d\n", query, userID)

    // 4. Execute the query
    rows, err := r.db.QueryContext(ctx, query, userID)
    if err != nil {
        return nil, fmt.Errorf("error querying todos: %w", err)
    }
    defer rows.Close()

    // 5. Initialize the slice with capacity
    todos := make([]domain.Todo, 0)

    // 6. Iterate over rows
    for rows.Next() {
        var todo domain.Todo
        err := rows.Scan(
            &todo.ID,
            &todo.UserID,
            &todo.Title,
            &todo.Description,
            &todo.Status,
            &todo.ImageID,
            &todo.CreatedAt,
            &todo.UpdatedAt,
        )
        if err != nil {
            return nil, fmt.Errorf("error scanning todo: %w", err)
        }
        todos = append(todos, todo)
    }

    // 7. Check for errors from iterating over rows
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("error iterating rows: %w", err)
    }

    return todos, nil
}
func (r *todoRepository) FindByID(ctx context.Context, id int) (*domain.Todo, error) {
	query := `SELECT id, user_id, title, description, status, image_id, created_at, updated_at 
              FROM todos 
              WHERE id = $1`

	var todo domain.Todo
	var imageID sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Description,
		&todo.Status,
		&imageID,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}

	if imageID.Valid {
		todo.ImageID = imageID.String
	}

	return &todo, nil
}

func (r *todoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	query := `INSERT INTO todos (user_id, title, description, status, image_id, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) 
              RETURNING id`

	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	err := r.db.QueryRowContext(
		ctx,
		query,
		todo.UserID,
		todo.Title,
		todo.Description,
		todo.Status,
		todo.ImageID,
		todo.CreatedAt,
		todo.UpdatedAt,
	).Scan(&todo.ID)

	return err
}

func (r *todoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	query := `UPDATE todos 
              SET title = $1, description = $2, status = $3, image_id = $4, updated_at = $5 
              WHERE id = $6`

	todo.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(
		ctx,
		query,
		todo.Title,
		todo.Description,
		todo.Status,
		todo.ImageID,
		todo.UpdatedAt,
		todo.ID,
	)

	return err
}


func (r *todoRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM todos WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	// Check if any row was affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("todo not found")
	}
	
	return nil
}