package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"

	"github.com/dendianugerah/bcke/internal/auth"
	"github.com/dendianugerah/bcke/internal/common/database"
	"github.com/dendianugerah/bcke/internal/common/middleware"
	"github.com/dendianugerah/bcke/internal/config"
	"github.com/dendianugerah/bcke/internal/user"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/dendianugerah/bcke/docs"
)

// @title           Core Backend API
// @version         1.0
// @description     A core backend API service with authentication and user management.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Dendi Anugerah
// @contact.email  dendianugrah40@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	log.Println("Starting server...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	// Connect to MongoDB
	client, err := database.ConnectMongoDB(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	db := client.Database(cfg.DBName)
	userCollection := db.Collection("users")

	// Initialize services and handlers
	userService := user.NewService(userCollection)
	userHandler := user.NewHandler(userService)

	authService := auth.NewService(userCollection, cfg.JWTSecret)
	authHandler := auth.NewHandler(authService)

	// Setup router with middleware
	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoveryMiddleware)

	// Public routes
	r.HandleFunc("/api/health", healthCheck).Methods(http.MethodGet)
	r.HandleFunc("/api/login", authHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/api/register", userHandler.Create).Methods(http.MethodPost)

	// Protected routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	api.HandleFunc("/users", userHandler.List).Methods(http.MethodGet)
	api.HandleFunc("/users/{id}", userHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/users/{id}", userHandler.Delete).Methods(http.MethodDelete)

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler())

	// Create server with timeouts
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server is running on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy"}`))
} 