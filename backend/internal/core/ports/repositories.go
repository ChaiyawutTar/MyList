package ports

import (

	"context"
	// "mime/multipart"

	"github.com/ChaiyawutTar/MyList/backend/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User, password string) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id int) (*domain.User, error)
}

