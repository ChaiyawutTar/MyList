package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"mime/multipart"
	"github.com/go-chi/chi/v5"
	"github.com/ChaiyawutTar/MyList/internal/core/domain"
	"github.com/ChaiyawutTar/MyList/internal/core/ports"
	"github.com/ChaiyawutTar/MyList/internal/adapters/handlers/middleware"
	"fmt"
	"time"
)

type TodoHandler struct {
	todoService ports.TodoService
}

func NewTodoHandler(todoService ports.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// internal/adapters/handlers/http/todo_handler.go
func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
    // Get user ID from context
    userID := middleware.GetUserIDFromContext(r.Context())

    // Check if we should include image data
    includeImages := r.URL.Query().Get("include_images") == "true"

    todos, err := h.todoService.GetAllTodos(r.Context(), userID)
    if err != nil {
        fmt.Printf("Error fetching todos: %v\n", err)
        http.Error(w, fmt.Sprintf("Failed to fetch todos: %v", err), http.StatusInternalServerError)
        return
    }

    // If we're not including images, remove the image data to reduce payload size
    if !includeImages {
        // Create a lightweight version of todos with just image IDs
        type LightTodo struct {
            ID          int       `json:"id"`
            UserID      int       `json:"user_id"`
            Title       string    `json:"title"`
            Description string    `json:"description"`
            Status      string    `json:"status"`
            ImageID     string    `json:"image_id,omitempty"`
            CreatedAt   time.Time `json:"created_at"`
            UpdatedAt   time.Time `json:"updated_at"`
        }

        lightTodos := make([]LightTodo, len(todos))
        for i, todo := range todos {
            lightTodos[i] = LightTodo{
                ID:          todo.ID,
                UserID:      todo.UserID,
                Title:       todo.Title,
                Description: todo.Description,
                Status:      todo.Status,
                ImageID:     todo.ImageID,
                CreatedAt:   todo.CreatedAt,
                UpdatedAt:   todo.UpdatedAt,
            }
        }

        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(lightTodos); err != nil {
            fmt.Printf("Error encoding todos: %v\n", err)
            http.Error(w, "Error encoding response", http.StatusInternalServerError)
        }
        return
    }

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(todos); err != nil {
        fmt.Printf("Error encoding todos: %v\n", err)
        http.Error(w, "Error encoding response", http.StatusInternalServerError)
    }
}
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
    // Get user ID from context
    userID := middleware.GetUserIDFromContext(r.Context())
    
    var req domain.CreateTodoRequest
    var imageFile *multipart.FileHeader
    
    // Parse multipart form for form-data (limit to 10MB)
    err := r.ParseMultipartForm(10 << 20)
    if err != nil {
        // If not multipart form, try JSON
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request format", http.StatusBadRequest)
            return
        }
    } else {
        // Get form values
        req = domain.CreateTodoRequest{
            Title:       r.FormValue("title"),
            Description: r.FormValue("description"),
            Status:      r.FormValue("status"),
        }
        
        // Get image file if present
        if file, header, err := r.FormFile("image"); err == nil {
            defer file.Close()
            imageFile = header
        }
    }
    
    // Validate required fields
    if req.Title == "" || req.Status == "" {
        http.Error(w, "Title and status are required", http.StatusBadRequest)
        return
    }

    // Create todo with image if provided
    todo, err := h.todoService.CreateTodo(r.Context(), req, userID, imageFile)
    if err != nil {
        fmt.Printf("Error creating todo: %v\n", err)
        http.Error(w, fmt.Sprintf("Failed to create todo: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(todo)
}

func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
    // Get todo ID from URL
    todoID, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        http.Error(w, "Invalid todo ID", http.StatusBadRequest)
        return
    }

    // Get user ID from context
    userID := middleware.GetUserIDFromContext(r.Context())
    
    var req domain.UpdateTodoRequest
    var imageFile *multipart.FileHeader
    
    // Parse multipart form for form-data (limit to 10MB)
    err = r.ParseMultipartForm(10 << 20)
    if err != nil {
        // If not multipart form, try JSON
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "Invalid request format", http.StatusBadRequest)
            return
        }
    } else {
        // Get form values
        req = domain.UpdateTodoRequest{
            Title:       r.FormValue("title"),
            Description: r.FormValue("description"),
            Status:      r.FormValue("status"),
        }
        
        // Get image file if present
        if file, header, err := r.FormFile("image"); err == nil {
            defer file.Close()
            imageFile = header
        }
    }
    
    // Validate required fields
    if req.Title == "" || req.Status == "" {
        http.Error(w, "Title and status are required", http.StatusBadRequest)
        return
    }

    // Update todo with image if provided
    todo, err := h.todoService.UpdateTodo(r.Context(), todoID, req, userID, imageFile)
    if err != nil {
        fmt.Printf("Error updating todo: %v\n", err)
        http.Error(w, fmt.Sprintf("Failed to update todo: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	// Get todo ID from URL
	todoID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	// Get user ID from context
	userID := middleware.GetUserIDFromContext(r.Context())

	// Delete todo
	if err := h.todoService.DeleteTodo(r.Context(), todoID, userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// internal/adapters/handlers/http/todo_handler.go
func (h *TodoHandler) GetTodoByID(w http.ResponseWriter, r *http.Request) {
    // Get todo ID from URL
    todoID, err := strconv.Atoi(chi.URLParam(r, "id"))
    if err != nil {
        http.Error(w, "Invalid todo ID", http.StatusBadRequest)
        return
    }

    // Get user ID from context
    userID := middleware.GetUserIDFromContext(r.Context())

    // Get todo
    todo, err := h.todoService.GetTodoByID(r.Context(), todoID, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(todo)
}