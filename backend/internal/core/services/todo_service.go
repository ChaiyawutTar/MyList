package services

import (
"context"
	"errors"
	"fmt"
	// "io"
	"mime/multipart"
	// "path/filepath"
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
    // Create a new todo without the image first
    todo := &domain.Todo{
        UserID:      userID,
        Title:       req.Title,
        Description: req.Description,
        Status:      req.Status,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    // If there's an image file, process it first before saving the todo
    if imageFile != nil {
        file, err := imageFile.Open()
        if err != nil {
            return nil, fmt.Errorf("failed to open image file: %w", err)
        }
        defer file.Close()
        
        // Save the image and get its ID
        imageID, err := s.imageRepo.Save(ctx, file, imageFile.Filename)
        if err != nil {
            return nil, fmt.Errorf("failed to save image: %w", err)
        }
        
        // Set the image ID on the todo
        todo.ImageID = imageID
    }
    
    // Now save the todo with the image ID
    err := s.todoRepo.Create(ctx, todo)  // Changed from createdTodo, err := to err :=
    if err != nil {
        // If we saved an image but failed to create the todo, clean up the image
        if todo.ImageID != "" {
            _ = s.imageRepo.Delete(ctx, todo.ImageID) // Best effort cleanup
        }
        return nil, fmt.Errorf("failed to create todo: %w", err)
    }
    
    return todo, nil  // Return the todo we created, not createdTodo
}

func (s *todoService) UpdateTodo(ctx context.Context, todoID int, req domain.UpdateTodoRequest, userID int, imageFile *multipart.FileHeader) (*domain.Todo, error) {
    // First, get the existing todo
    existingTodo, err := s.todoRepo.FindByID(ctx, todoID)
    if err != nil {
        return nil, fmt.Errorf("failed to find todo: %w", err)
    }
    
    // Check if the todo belongs to the user
    if existingTodo.UserID != userID {
        return nil, errors.New("unauthorized")
    }
    
    // Update the todo fields
    existingTodo.Title = req.Title
    existingTodo.Description = req.Description
    existingTodo.Status = req.Status
    existingTodo.UpdatedAt = time.Now()
    
    // If there's a new image file, process it
    if imageFile != nil {
        file, err := imageFile.Open()
        if err != nil {
            return nil, fmt.Errorf("failed to open image file: %w", err)
        }
        defer file.Close()
        
        // Save the new image
        imageID, err := s.imageRepo.Save(ctx, file, imageFile.Filename)
        if err != nil {
            return nil, fmt.Errorf("failed to save image: %w", err)
        }
        
        // If there was an existing image, delete it
        if existingTodo.ImageID != "" {
            _ = s.imageRepo.Delete(ctx, existingTodo.ImageID) // Best effort cleanup
        }
        
        // Set the new image ID
        existingTodo.ImageID = imageID
    }
    
    // Update the todo in the database
    err = s.todoRepo.Update(ctx, existingTodo)  // Changed from updatedTodo, err := to err :=
    if err != nil {
        // If we saved a new image but failed to update the todo, clean up the new image
        if imageFile != nil && existingTodo.ImageID != "" {
            _ = s.imageRepo.Delete(ctx, existingTodo.ImageID) // Best effort cleanup
        }
        return nil, fmt.Errorf("failed to update todo: %w", err)
    }
    
    return existingTodo, nil  // Return the todo we updated, not updatedTodo
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