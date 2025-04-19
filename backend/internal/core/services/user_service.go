package services

import (
	"context"
	"errors"
	// "fmt"
	"time"
	"github.com/ChaiyawutTar/MyList/internal/core/domain"
	"github.com/ChaiyawutTar/MyList/internal/core/ports"
	"github.com/ChaiyawutTar/MyList/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo ports.UserRepository
	jwtAuth  *auth.JWTAuth
}

func NewUserService(userRepo ports.UserRepository, jwtAuth *auth.JWTAuth) ports.UserService {
	return &userService{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
	}
}

func (s *userService) CreateUser(ctx context.Context, req domain.SignupRequest) (*domain.AuthResponse, error) {
	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, errors.New("username, email and password are required")
	}

	// Create user
	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user, req.Password); err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.jwtAuth.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *userService) Login(ctx context.Context, req domain.LoginRequest) (*domain.AuthResponse, error) {
	// Validate input
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("email and password are required")
	}

	// Find user
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate token
	token, err := s.jwtAuth.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *userService) OAuthLogin(ctx context.Context, oauthUser domain.OAuthUser) (*domain.AuthResponse, error) {
    // Try to find user by OAuth provider ID
    user, err := s.userRepo.FindByOAuthID(ctx, oauthUser.Provider, oauthUser.ProviderID)
    
    if err != nil {
        // Try to find user by email
        user, err = s.userRepo.FindByEmail(ctx, oauthUser.Email)
        
        if err != nil {
            // Create new user
            user = &domain.User{
                Username:        oauthUser.Name,
                Email:           oauthUser.Email,
                OAuthProvider:   oauthUser.Provider,
                OAuthProviderID: oauthUser.ProviderID,
                CreatedAt:       time.Now(),
            }
            
            // Create user without password
            if err := s.userRepo.Create(ctx, user, ""); err != nil {
                return nil, err
            }
        } else {
            // Update existing user with OAuth info
            user.OAuthProvider = oauthUser.Provider
            user.OAuthProviderID = oauthUser.ProviderID
            
            // Update user in database
            if err := s.userRepo.UpdateOAuthInfo(ctx, user); err != nil {
                return nil, err
            }
        }
    }
    
    // Generate token
    token, err := s.jwtAuth.GenerateToken(user.ID)
    if err != nil {
        return nil, err
    }
    
    return &domain.AuthResponse{
        Token: token,
        User:  *user,
    }, nil
}