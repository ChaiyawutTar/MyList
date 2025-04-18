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
)

type TodoHandler struct {
	todoService ports.TodoService
}

func NewTodoHandler(todoService ports.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

func (h *TodoHandler) GetAllTodos(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := middleware.GetUserIDFromContext(r.Context())

	todos, err := h.todoService.GetAllTodos(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID := middleware.GetUserIDFromContext(r.Context())
	
	var req domain.CreateTodoRequest
	var imageFile *multipart.FileHeader
	
	// Check content type to determine how to parse the request
	contentType := r.Header.Get("Content-Type")
	
	if contentType == "application/json" {
		// Parse JSON request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		// Parse multipart form for form-data
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			// If not multipart form, try regular form
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}
		}
		
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

	// Create todo
	todo, err := h.todoService.CreateTodo(r.Context(), req, userID, imageFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	
	// Check content type to determine how to parse the request
	contentType := r.Header.Get("Content-Type")
	
	if contentType == "application/json" {
		// Parse JSON request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
	} else {
		// Parse multipart form for form-data
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			// If not multipart form, try regular form
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Error parsing form", http.StatusBadRequest)
				return
			}
		}
		
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

	// Update todo
	todo, err := h.todoService.UpdateTodo(r.Context(), todoID, req, userID, imageFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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