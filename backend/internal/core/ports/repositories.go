package ports

import (

	"context"
	"mime/multipart"

	"github.com/ChaiyawutTar/MyList/internal/core/domain"
)

type UserRepository interface {
    Create(ctx context.Context, user *domain.User, password string) error
    FindByEmail(ctx context.Context, email string) (*domain.User, error)
    FindByID(ctx context.Context, id int) (*domain.User, error)
    FindByOAuthID(ctx context.Context, provider, providerID string) (*domain.User, error)
    UpdateOAuthInfo(ctx context.Context, user *domain.User) error
}

type TodoRepository interface {
	FindAll(ctx context.Context, userID int) ([]domain.Todo, error)
	FindByID(ctx context.Context, id int) (*domain.Todo, error)
	Create(ctx context.Context, todo *domain.Todo) error
	Update(ctx context.Context, todo *domain.Todo) error
	Delete(ctx context.Context, id int) error
}

type ImageRepository interface {
	Save(ctx context.Context, file multipart.File, filename string) (string, error)
	Delete(ctx context.Context, path string) error
}