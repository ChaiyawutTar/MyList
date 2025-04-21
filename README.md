# MyList - Task Management Application

MyList is a modern task management application built with Go (backend) and Next.js (frontend). It features user authentication (including OAuth), todo management with image attachments, and a clean, responsive UI built with Shadcn UI components.

## 🚀 Features

- **User Authentication**
  - Traditional email/password signup and login
  - Google OAuth integration
  - JWT-based authentication

- **Todo Management**
  - Create, read, update, and delete todos
  - Attach images to todos
  - Filter todos by status

- **Modern UI**
  - Responsive design using Shadcn UI components
  - Optimized image loading

## 🏗️ Architecture

The application follows a clean architecture approach with clear separation of concerns:

### Backend (Go)

```
├── internal
│   ├── adapters        # Implementation of interfaces
│   │   ├── handlers    # HTTP handlers
│   │   └── repositories # Data access layer
│   ├── config          # Application configuration
│   └── core            # Business logic
│       ├── domain      # Domain models
│       ├── ports       # Interface definitions
│       └── services    # Business logic implementation
├── pkg                 # Reusable packages
│   └── auth            # Authentication utilities
└── main.go             # Application entry point
```

### Frontend (Next.js)

The frontend is built with Next.js and uses a clean architecture approach with:

- Components: Reusable UI components
- Hooks: Custom React hooks for state management
- Services: API communication layer
- Providers: Context providers for global state

## 🛠️ Technology Stack

### Backend
- **Language: [Go](https://go.dev/)**
- **Web Framework: [Chi router](https://go-chi.io/#/)**
- **Database: [PostgreSQL (hosted on Neon)](https://www.neon.tech)**
- **Authentication**: JWT + Google OAuth
- **Image Storage: [PostgreSQL (hosted on Neon) store in byte file](https://www.neon.tech)**

### Frontend
- **Framework: [Next.js](https://nextjs.org/)**
- **UI Library: [Shadcn UI](https://ui.shadcn.com/)**
- **HTTP Client: [Axios](https://axios-http.com/)**

## 📊 Database Schema

The application uses a PostgreSQL database with the following main tables:

- **users**: Stores user information and authentication details
- **todos**: Stores todo items with references to users
- **images**: Stores image data for todo attachments
```sql
CREATE TABLE users (

    id SERIAL PRIMARY KEY,

    username VARCHAR(50) UNIQUE NOT NULL,

    email VARCHAR(100) UNIQUE NOT NULL,

    password_hash VARCHAR(255) NOT NULL,

oauth_provider VARCHAR(50),

oauth_provider_id VARCHAR(255),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE todos (

    id SERIAL PRIMARY KEY,

    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,

    title VARCHAR(100) NOT NULL,

    description TEXT,

    status VARCHAR(20) DEFAULT 'pending',

    image_id VARCHAR(255),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

/* for the image table */

  CREATE TABLE images (

    id SERIAL PRIMARY KEY,

    filename TEXT NOT NULL,

    data BYTEA NOT NULL,

    content_type TEXT,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);
```

## 🚀 Deployment

### Backend
The backend is deployed on [Render](https://render.com) at:
```
https://mylist-deployment.onrender.com
```

### Frontend
The frontend is deployed on [Vercel](https://vercel.com) at:
```
https://mylist-delta.vercel.app
```

## 🏁 Getting Started

### Prerequisites
- Go 1.20+
- Node.js 18+
- PostgreSQL database

### Local Development

#### Backend
1. Clone the repository
2. Set up environment variables (see `.env.example`)
3. Run the application:
```bash
cd backend
go run main.go
```

#### Frontend
1. Navigate to the frontend directory
2. Install dependencies:
```bash
cd frontend
npm install
```
3. Start the development server:
```bash
npm run dev
```

## 🧩 Code Explanation

### Clean Architecture (Hexagonal Architecture)

The project follows clean architecture principles to ensure:

- **Separation of Concerns**: Each layer has a specific responsibility
- **Dependency Inversion**: High-level modules don't depend on low-level modules
- **Testability**: Business logic can be tested independently of external dependencies

### Key Design Decisions

1. **Repository Pattern**: Abstracts data access logic, making it easy to switch between different data sources
2. **Dependency Injection**: Services and handlers receive their dependencies through constructors
3. **Interface-Based Design**: Components interact through interfaces, not concrete implementations
4. **Database-Based Image Storage**: Images are stored directly in the database for simplicity and portability

### Authentication Flow

1. **Traditional Auth**: Users can register and login with email/password
2. **OAuth Flow**: 
   - User initiates OAuth flow by clicking "Sign in with Google"
   - Backend redirects to Google for authentication
   - Google redirects back to the callback URL with an authorization code
   - Backend exchanges the code for tokens and user information
   - User is authenticated and redirected to the application

### Todo Management

Todos are managed through a RESTful API with endpoints for:
- Creating new todos (with optional image attachments)
- Retrieving todos (individual or all)
- Updating existing todos
- Deleting todos

## 🔷 Hexagonal Architecture (Ports and Adapters)

The MyList application implements hexagonal architecture (also known as ports and adapters), a pattern that puts the domain logic at the center of the application and defines how it interacts with the outside world through ports and adapters.

### Core Principles

1. **Domain-Centric Design**: The business logic is at the center, isolated from external concerns
2. **Ports**: Interfaces that define how the domain interacts with the outside world
3. **Adapters**: Implementations of these interfaces that connect to specific technologies

### Implementation in MyList

#### Domain Layer (`internal/core/domain`)

The domain layer contains the business entities and logic, completely isolated from external concerns:

- `todo.go`: Defines the Todo entity with its properties and behaviors
- `user.go`: Defines the User entity with authentication-related behaviors

These domain models are pure Go structs with no dependencies on external frameworks or libraries.

#### Ports Layer (`internal/core/ports`)

The ports layer defines interfaces that the domain uses to interact with the outside world:

- `repositories.go`: Defines interfaces for data persistence operations
  
  ```go
  type TodoRepository interface {
      Create(todo *domain.Todo) (*domain.Todo, error)
      GetByID(id string) (*domain.Todo, error)
      GetAllByUserID(userID string) ([]*domain.Todo, error)
      Update(todo *domain.Todo) (*domain.Todo, error)
      Delete(id string) error
  }
  
  type UserRepository interface {
      Create(user *domain.User) (*domain.User, error)
      GetByID(id string) (*domain.User, error)
      GetByEmail(email string) (*domain.User, error)
      GetByOAuthID(provider, oauthID string) (*domain.User, error)
      Update(user *domain.User) (*domain.User, error)
  }
  
  type ImageRepository interface {
      Save(data []byte, contentType string) (string, error)
      Get(id string) ([]byte, string, error)
      Delete(id string) error
  }
  ```
  
#### Service Layer (`internal/core/services`)
The service layer implements the application's use cases by orchestrating the domain entities and interacting with ports:

- `todo_service.go`: Implements the TodoService interface
- `user_service.go`: Implements the UserService interface
- `services.go`: Defines interfaces for application services

These services contain the application's business logic but remain agnostic to how data is stored or how the UI is implemented.

  ```go
  type TodoService interface {
      CreateTodo(userID, title, description string, imageData []byte, imageType string) (*domain.Todo, error)
      GetTodoByID(id string) (*domain.Todo, error)
      GetAllTodosByUserID(userID string) ([]*domain.Todo, error)
      UpdateTodo(id, title, description, status string, imageData []byte, imageType string) (*domain.Todo, error)
      DeleteTodo(id string) error
  }
  
  type UserService interface {
      Signup(username, email, password string) (*domain.User, string, error)
      Login(email, password string) (*domain.User, string, error)
      GetUserByID(id string) (*domain.User, error)
  }
  ```
  

### Adapters Layer (`internal/adapters`)
The adapters layer contains implementations of the ports that connect to specific technologies:

- Repository Adapters (`internal/adapters/repositories/postgres`):
    - `todo_repository.go`: PostgreSQL implementation of TodoRepository
    - `user_repository.go`: PostgreSQL implementation of UserRepository
    - `image_repository.go`: PostgreSQL implementation of ImageRepository
- HTTP Adapters (`internal/adapters/handlers/http`):
    - `todo_handler.go`: HTTP handlers for todo operations
    - `auth_handler.go`: HTTP handlers for authentication
    - `image_handler.go`: HTTP handlers for image operations
## Infrastructure of Hexagonal Architecture
![image](https://github.com/user-attachments/assets/992aa05f-467c-44df-87a1-0ab5d90fee67)

## 📝 API Documentation

### Authentication Endpoints

- `POST /signup`: Register a new user
- `POST /login`: Authenticate a user
- `GET /auth/{provider}`: Initiate OAuth flow
- `GET /auth/{provider}/callback`: Handle OAuth callback

### Todo Endpoints

- `GET /todos`: Get all todos for the authenticated user
- `POST /todos`: Create a new todo
- `GET /todos/{id}`: Get a specific todo
- `PUT /todos/{id}`: Update a todo
- `DELETE /todos/{id}`: Delete a todo

### Image Endpoints

- `GET /images/{id}`: Retrieve an image by ID

## 🔒 Security Considerations

- JWT tokens are used for authentication
- Passwords are hashed before storage
- CORS is configured to allow only specific origins
- Input validation is performed on all endpoints

