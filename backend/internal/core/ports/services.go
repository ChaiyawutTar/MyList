package ports

import (

	"context"
	"mime/multipart"

	"github.com/ChaiyawutTar/MyList/internal/core/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, req domain.SignupRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
}

type TodoService interface {
	GetAllTodos(ctx context.Context, userID int) ([]domain.Todo, error)
	GetTodoByID(ctx context.Context, id int, userID int) (*domain.Todo, error)
	CreateTodo(ctx context.Context, req domain.CreateTodoRequest, userID int, image *multipart.FileHeader) (*domain.Todo, error)
	UpdateTodo(ctx context.Context, id int, req domain.UpdateTodoRequest, userID int, image *multipart.FileHeader) (*domain.Todo, error)
	DeleteTodo(ctx context.Context, id int, userID int) error
}