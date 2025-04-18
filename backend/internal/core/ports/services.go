package ports

import (

	"context"
	// "mime/multipart"

	"github.com/ChaiyawutTar/MyList/backend/internal/core/domain"
)

type UserService interface {
	CreateUser(ctx context.Context, req domain.SignupRequest) (*domain.AuthResponse, error)
	Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
}