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
  - Dark/light mode support
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
- **Language**: Go
- **Web Framework**: Chi router
- **Database**: PostgreSQL (hosted on Neon)
- **Authentication**: JWT + Google OAuth
- **Image Storage**: Database-based storage

### Frontend
- **Framework**: Next.js
- **UI Library**: Shadcn UI
- **State Management**: React Context
- **HTTP Client**: Axios

## 📊 Database Schema

The application uses a PostgreSQL database with the following main tables:

- **users**: Stores user information and authentication details
- **todos**: Stores todo items with references to users
- **images**: Stores image data for todo attachments

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
go run main.go
```

#### Frontend
1. Navigate to the frontend directory
2. Install dependencies:
```bash
npm install
```
3. Start the development server:
```bash
npm run dev
```

## 🧩 Code Explanation

### Clean Architecture

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

## 🧪 Testing

The application includes unit tests for core business logic and integration tests for API endpoints.

To run tests:

```bash
go test ./...
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## Contact

For any questions or feedback, please reach out to [your-email@example.com](mailto:your-email@example.com).
