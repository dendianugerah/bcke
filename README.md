# Core Backend API

A robust and scalable backend API service built with Go

## 📋 Prerequisites

Before you begin, ensure you have the following installed:
- Go 1.22 or higher
- MongoDB 4.4 or higher
- Docker and Docker Compose (optional)

## 🛠️ Installation

### Local Development

1. Clone the repository
```bash
git clone https://github.com/dendianugerah/bcke.git
cd bcke
```

2. Install dependencies
```bash
go mod download
```

3. Set up environment variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the application
```bash
# Direct execution
go run cmd/api/main.go

# Or using Docker
docker-compose up --build
```

### Docker Deployment

Run the entire stack using Docker Compose:
```bash
docker-compose up --build
```

This will start:
- The API server on port 8080
- MongoDB instance on port 27017

## 🔑 API Authentication

The API uses JWT (JSON Web Tokens) for authentication. To access protected endpoints:

1. Register a new user
2. Login to get a JWT token
3. Include the token in subsequent requests:
   ```
   Authorization: Bearer <your_token>
   ```

## 📚 API Documentation

Access the Swagger documentation at: `http://localhost:8080/swagger/`

## 🏗️ Project Structure

```
.
├── cmd/
│   └── api/            # Application entrypoint
│       └── main.go
├── internal/
│   ├── auth/           # Authentication package
│   │   ├── handler.go
│   │   ├── model.go
│   │   └── service.go
│   ├── common/         # Shared components
│   │   ├── database/
│   │   ├── middleware/
│   │   └── response/
│   ├── config/         # Configuration
│   └── user/           # User management
│       ├── handler.go
│       ├── model.go
│       └── service.go
├── docs/               # Swagger documentation
├── docker-compose.yml
├── Dockerfile
└── README.md
```

## ⚙️ Configuration

Environment variables (`.env`):
```env
MONGODB_URI=mongodb://localhost:27017
DB_NAME=core_backend
JWT_SECRET=your-secret-key
PORT=8080
```

### Generating Documentation
```bash
# Generate Swagger docs
swag init -g cmd/api/main.go
```

## 📧 Contact

Dendi Anugerah - dendianugrah40@gmail.com

Project Link: [https://github.com/dendianugerah/bcke](https://github.com/dendianugerah/bcke)