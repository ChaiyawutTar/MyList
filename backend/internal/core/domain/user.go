// internal/core/domain/user.go
package domain

import "time"

type User struct {
    ID              int       `json:"id"`
    Username        string    `json:"username"`
    Email           string    `json:"email"`
    PasswordHash    string    `json:"-"`
    OAuthProvider   string    `json:"oauth_provider,omitempty"`
    OAuthProviderID string    `json:"oauth_provider_id,omitempty"`
    CreatedAt       time.Time `json:"created_at"`
}

// Add OAuthUser struct
type OAuthUser struct {
    Provider    string `json:"provider"`
    Email       string `json:"email"`
    Name        string `json:"name"`
    AvatarURL   string `json:"avatar_url,omitempty"`
    ProviderID  string `json:"provider_id"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
