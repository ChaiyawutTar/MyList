package postgres

import (
	"database/sql"
	"context"
	"errors"
	"time"
	"github.com/ChaiyawutTar/MyList/internal/core/domain"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User, password string) error {
    // Hash password if provided
    var hashedPassword string
    if password != "" {
        hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        hashedPassword = string(hashed)
    }

    // Insert user into database
    query := `INSERT INTO users (username, email, password_hash, oauth_provider, oauth_provider_id, created_at)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id`

    err := r.db.QueryRowContext(
        ctx,
        query,
        user.Username,
        user.Email,
        hashedPassword,
        user.OAuthProvider,
        user.OAuthProviderID,
        time.Now(),
    ).Scan(&user.ID)

    return err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, created_at 
              FROM users 
              WHERE email = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id int) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, created_at 
              FROM users 
              WHERE id = $1`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindByOAuthID(ctx context.Context, provider, providerID string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, created_at
	FROM users
	WHERE oauth_provider = $1 AND oauth_provider_id = $2`

	var user domain.User
	err := r.db.QueryRowContext(ctx, query, provider, providerID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateOAuthInfo(ctx context.Context, user *domain.User) error {
    query := `UPDATE users 
              SET oauth_provider = $1, oauth_provider_id = $2
              WHERE id = $3`

    _, err := r.db.ExecContext(
        ctx,
        query,
        user.OAuthProvider,
        user.OAuthProviderID,
        user.ID,
    )

    return err
}