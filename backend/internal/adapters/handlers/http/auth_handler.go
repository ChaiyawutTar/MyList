package http

import (
    "encoding/json"
    "fmt"
    "net/http"
	"log"
    
    "github.com/markbates/goth/gothic"
    "github.com/ChaiyawutTar/MyList/internal/core/domain"
    "github.com/ChaiyawutTar/MyList/internal/core/ports"
    "github.com/ChaiyawutTar/MyList/pkg/auth"
    "github.com/go-chi/chi/v5"
)

type AuthHandler struct {
    userService ports.UserService
    oauthManager *auth.OAuthManager
    frontendURL string
}

func NewAuthHandler(userService ports.UserService, oauthManager *auth.OAuthManager, frontendURL string) *AuthHandler {
    return &AuthHandler{
        userService: userService,
        oauthManager: oauthManager,
        frontendURL: frontendURL,
    }
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
    var req domain.SignupRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    resp, err := h.userService.CreateUser(r.Context(), req)
    if err != nil {
        // Log the error for debugging
        log.Printf("Signup error: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req domain.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    resp, err := h.userService.Login(r.Context(), req)
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) BeginOAuth(w http.ResponseWriter, r *http.Request) {
    provider := chi.URLParam(r, "provider")
    gothic.GetProviderName = func(*http.Request) (string, error) {
        return provider, nil
    }
    
    // Set the callback URL in the session
    gothic.BeginAuthHandler(w, r)
}

func (h *AuthHandler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
    provider := chi.URLParam(r, "provider")
    gothic.GetProviderName = func(*http.Request) (string, error) {
        return provider, nil
    }
    
    // Complete the auth process
    user, err := gothic.CompleteUserAuth(w, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Process the user data and create/login the user
    resp, err := h.userService.OAuthLogin(r.Context(), domain.OAuthUser{
        Provider:    user.Provider,
        Email:       user.Email,
        Name:        user.Name,
        AvatarURL:   user.AvatarURL,
        ProviderID:  user.UserID,
    })
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    redirectURL := fmt.Sprintf("%s/callback?token=%s", h.frontendURL, resp.Token)
    http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}