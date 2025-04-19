package services

import (
"context"
	"errors"
	"fmt"
	// "io"
	"mime/multipart"
	"path/filepath"
	"time"
	"github.com/ChaiyawutTar/MyList/internal/core/domain"
	"github.com/ChaiyawutTar/MyList/internal/core/ports"
)

type todoService struct {
	todoRepo  ports.TodoRepository
	imageRepo ports.ImageRepository
}

func NewTodoService(todoRepo ports.TodoRepository, imageRepo ports.ImageRepository) ports.TodoService {
	return &todoService{
		todoRepo:  todoRepo,
		imageRepo: imageRepo,
	}
}

func (s *todoService) GetAllTodos(ctx context.Context, userID int) ([]domain.Todo, error) {
    todos, err := s.todoRepo.FindAll(ctx, userID)
    if err != nil {
        fmt.Printf("Error in todoService.GetAllTodos: %v\n", err)
        return nil, err
    }
    return todos, nil
}

func (s *todoService) GetTodoByID(ctx context.Context, id int, userID int) (*domain.Todo, error) {
	todo, err := s.todoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if todo.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return todo, nil
}

func (s *todoService) CreateTodo(ctx context.Context, req domain.CreateTodoRequest, userID int, imageFile *multipart.FileHeader) (*domain.Todo, error) {
	// Validate input
	if req.Title == "" {
		return nil, errors.New("title is required")
	}

	// Set default status if not provided
	if req.Status == "" {
		req.Status = "pending"
	}

	// Create todo
	todo := &domain.Todo{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Handle image upload
	if imageFile != nil {
		file, err := imageFile.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// Generate unique filename
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(imageFile.Filename))
		
		// Save image to database
		imageID, err := s.imageRepo.Save(ctx, file, filename)
		if err != nil {
			return nil, err
		}

		todo.ImageID = imageID // Changed from ImagePath
	}

	// Save todo
	if err := s.todoRepo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) UpdateTodo(ctx context.Context, id int, req domain.UpdateTodoRequest, userID int, imageFile *multipart.FileHeader) (*domain.Todo, error) {
	// Get existing todo
	todo, err := s.todoRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if todo.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Update fields
	todo.Title = req.Title
	todo.Description = req.Description
	todo.Status = req.Status
	todo.UpdatedAt = time.Now()

	// Handle image upload
	if imageFile != nil {
		file, err := imageFile.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// Delete old image if exists
		if todo.ImageID != "" {
			if err := s.imageRepo.Delete(ctx, todo.ImageID); err != nil {
				// Log error but continue
				fmt.Printf("Error deleting old image: %v\n", err)
			}
		}

		// Generate unique filename
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(imageFile.Filename))
		
		// Save new image
		imageID, err := s.imageRepo.Save(ctx, file, filename)
		if err != nil {
			return nil, err
		}

		todo.ImageID = imageID // Changed from ImagePath
	}

	// Save todo
	if err := s.todoRepo.Update(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) DeleteTodo(ctx context.Context, id int, userID int) error {
	// Get existing todo
	todo, err := s.todoRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Verify ownership
	if todo.UserID != userID {
		return errors.New("unauthorized")
	}

	// Delete image if exists
	if todo.ImageID != "" {
		if err := s.imageRepo.Delete(ctx, todo.ImageID); err != nil {
			// Log error but continue
			fmt.Printf("Error deleting image: %v\n", err)
		}
	}

	// Delete todo
	return s.todoRepo.Delete(ctx, id)
}