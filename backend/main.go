package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	// "os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth/gothic"

	httphandlers "github.com/ChaiyawutTar/MyList/internal/adapters/handlers/http"
	custommiddleware "github.com/ChaiyawutTar/MyList/internal/adapters/handlers/middleware"

	// Remove or comment out the file repository import
	// "github.com/ChaiyawutTar/MyList/internal/adapters/repositories/file"
	"github.com/ChaiyawutTar/MyList/internal/adapters/repositories/postgres"
	// "github.com/ChaiyawutTar/MyList/internal/adapters/repositories/file"

	"github.com/ChaiyawutTar/MyList/internal/config"
	"github.com/ChaiyawutTar/MyList/internal/core/services"
	"github.com/ChaiyawutTar/MyList/pkg/auth"


)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Initialize JWT auth
	jwtAuth := auth.NewJWTAuth(cfg.JWTSecret, cfg.JWTExpiry)

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	todoRepo := postgres.NewTodoRepository(db)
	
	// Replace file-based image repository with database-based one
	imageRepo := postgres.NewImageRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo, jwtAuth)
	todoService := services.NewTodoService(todoRepo, imageRepo)

	// Initialize handlers
	todoHandler := httphandlers.NewTodoHandler(todoService)
	
	// Add image handler for serving images from database
	imageHandler := httphandlers.NewImageHandler(imageRepo)
	
	oauthManager := auth.NewOAuthManager(
		userRepo,
		jwtAuth,
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.OAuthCallbackURL,
	)
	authHandler := httphandlers.NewAuthHandler(userService, oauthManager, cfg.FrontendURL)

	sessionStore := sessions.NewCookieStore([]byte(cfg.SessionSecret))
	sessionStore.MaxAge(86400 * 30) // 30 days
	sessionStore.Options.Path = "/"
	sessionStore.Options.HttpOnly = true
	sessionStore.Options.Secure = false // Set to true in production with HTTPS
	
	gothic.Store = sessionStore

	// Initialize router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           300,
	}))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/signup", authHandler.Signup)
		r.Post("/login", authHandler.Login)

		r.Get("/auth/{provider}", authHandler.BeginOAuth)
		r.Get("/auth/{provider}/callback", authHandler.OAuthCallback)
		
		// Add route for serving images from database
		r.Get("/images/{id}", imageHandler.ServeImage)
	})

	// Remove or comment out the static file server since images are now in the database
	// r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.UploadDir))))

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(custommiddleware.AuthMiddleware(jwtAuth))

		r.Get("/todos", todoHandler.GetAllTodos)
		r.Post("/todos", todoHandler.CreateTodo)
		r.Get("/todos/{id}", todoHandler.GetTodoByID)
		r.Put("/todos/{id}", todoHandler.UpdateTodo)
		r.Delete("/todos/{id}", todoHandler.DeleteTodo)
	})

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	fmt.Printf("Server started on %s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}