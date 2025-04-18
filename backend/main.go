package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	// "time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"

	httphandlers "github.com/ChaiyawutTar/MyList/backend/internal/adapters/handlers/http"
	// custommiddleware "github.com/ChaiyawutTar/MyList/backend/internal/adapters/handlers/middleware"
	"github.com/ChaiyawutTar/MyList/backend/internal/config"
	"github.com/ChaiyawutTar/MyList/backend/pkg/auth"
	"github.com/ChaiyawutTar/MyList/backend/internal/core/services"
	"github.com/ChaiyawutTar/MyList/backend/internal/adapters/repositories/postgres"

)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DatabaseURL)

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatal(err)
	}
	// Initialize JWT auth
	jwtAuth := auth.NewJWTAuth(cfg.JWTSecret, cfg.JWTExpiry)

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)

	// Initialize service
	userService := services.NewUserService(userRepo, jwtAuth)

	authHandler := httphandlers.NewAuthHandler(userService)

	// Initialize router
	r := chi.NewRouter()
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

	r.Group(func(r chi.Router) {
		r.Post("/signup", authHandler.Signup)
		r.Post("/login", authHandler.Login)
	})

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	fmt.Printf("Server started on %s\n", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}