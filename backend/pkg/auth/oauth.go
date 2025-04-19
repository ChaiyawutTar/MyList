// pkg/auth/oauth.go
package auth

import (
	"context"
	// "encoding/json"
	// "errors"
	// "fmt"
	"net/http"
	"time"

	"github.com/ChaiyawutTar/MyList/internal/core/domain"
	"github.com/ChaiyawutTar/MyList/internal/core/ports"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

// OAuthManager handles OAuth authentication
type OAuthManager struct {
	userRepo ports.UserRepository
	jwtAuth  *JWTAuth
}

// NewOAuthManager creates a new OAuthManager
func NewOAuthManager(userRepo ports.UserRepository, jwtAuth *JWTAuth, googleClientID, googleClientSecret, callbackURL string) *OAuthManager {
	// Configure Goth
	goth.UseProviders(
		google.New(googleClientID, googleClientSecret, callbackURL, "email", "profile"),
	)

	return &OAuthManager{
		userRepo: userRepo,
		jwtAuth:  jwtAuth,
	}
}

// HandleCallback processes OAuth callback and returns user info
func (m *OAuthManager) HandleCallback(ctx context.Context, w http.ResponseWriter, r *http.Request, provider string) (*domain.AuthResponse, error) {
    // Set the provider name in the request context
    gothic.GetProviderName = func(*http.Request) (string, error) {
        return provider, nil
    }
    
    // Get the user from the callback
    gothUser, err := gothic.CompleteUserAuth(w, r)
    if err != nil {
        return nil, err
    }

    // Check if user exists by OAuth provider ID
    user, err := m.userRepo.FindByOAuthID(ctx, gothUser.Provider, gothUser.UserID)
    if err != nil {
        // Try to find user by email
        user, err = m.userRepo.FindByEmail(ctx, gothUser.Email)
        if err != nil {
            // User doesn't exist, create a new one
            newUser := &domain.User{
                Username:        gothUser.Name,
                Email:           gothUser.Email,
                OAuthProvider:   gothUser.Provider,
                OAuthProviderID: gothUser.UserID,
                CreatedAt:       time.Now(),
            }

            // Create user without password
            if err := m.userRepo.Create(ctx, newUser, ""); err != nil {
                return nil, err
            }
            
            user = newUser
        } else {
            // Update existing user with OAuth info
            user.OAuthProvider = gothUser.Provider
            user.OAuthProviderID = gothUser.UserID
            
            if err := m.userRepo.UpdateOAuthInfo(ctx, user); err != nil {
                return nil, err
            }
        }
    }

    // Generate JWT token
    token, err := m.jwtAuth.GenerateToken(user.ID)
    if err != nil {
        return nil, err
    }

    return &domain.AuthResponse{
        Token: token,
        User:  *user,
    }, nil
}

// func BeginAuthHandler(w http.ResponseWriter, r *http.Request, provider string) {
//     gothic.GetProviderName = func(*http.Request) (string, error) {
//         return provider, nil
//     }
//     gothic.BeginAuthHandler(w, r)
// }