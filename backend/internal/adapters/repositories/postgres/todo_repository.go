package postgres

import (
	"context"
	"database/sql"

	"errors"
	"time"

	"github.com/ChaiyawutTar/MyList/internal/core/domain"
)

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) *todoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) FindAll(ctx context.Context, userID int) ([]domain.Todo, error) {
	query := `SELECT id, user_id, title, description, status, image_path, created_at, updated_at 
              FROM todos 
              WHERE user_id = $1 
              ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []domain.Todo
	for rows.Next() {
		var todo domain.Todo
		var imagePath sql.NullString
		err := rows.Scan(
			&todo.ID,
			&todo.UserID,
			&todo.Title,
			&todo.Description,
			&todo.Status,
			&imagePath,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if imagePath.Valid {
			todo.ImagePath = imagePath.String
		}

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *todoRepository) FindByID(ctx context.Context, id int) (*domain.Todo, error) {
	query := `SELECT id, user_id, title, description, status, image_path, created_at, updated_at 
              FROM todos 
              WHERE id = $1`

	var todo domain.Todo
	var imagePath sql.NullString
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&todo.ID,
		&todo.UserID,
		&todo.Title,
		&todo.Description,
		&todo.Status,
		&imagePath,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("todo not found")
		}
		return nil, err
	}

	if imagePath.Valid {
		todo.ImagePath = imagePath.String
	}

	return &todo, nil
}

func (r *todoRepository) Create(ctx context.Context, todo *domain.Todo) error {
	query := `INSERT INTO todos (user_id, title, description, status, image_path, created_at, updated_at) 
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
		todo.ImagePath,
		todo.CreatedAt,
		todo.UpdatedAt,
	).Scan(&todo.ID)

	return err
}

func (r *todoRepository) Update(ctx context.Context, todo *domain.Todo) error {
	query := `UPDATE todos 
              SET title = $1, description = $2, status = $3, image_path = $4, updated_at = $5 
              WHERE id = $6`

	todo.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(
		ctx,
		query,
		todo.Title,
		todo.Description,
		todo.Status,
		todo.ImagePath,
		todo.UpdatedAt,
		todo.ID,
	)

	return err
}

func (r *todoRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM todos WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}