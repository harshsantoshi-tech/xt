// main.go
package main

import (
	"context"
	"errors"
	"expense-tracker/configs"
	"expense-tracker/routes"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// A global Echo instance to be reused across requests.
var e *echo.Echo

func init() {
	// Initialize the Echo server and routes once, when the application starts.
	e = echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register all the API routes
	routes.InitRoutes(e)
}

// Handler is the Vercel-compatible entry point for your serverless function.
// It simply serves the incoming request using the global Echo instance.
func Handler(w http.ResponseWriter, r *http.Request) {
	e.ServeHTTP(w, r)
}

func main() {
	// Start the server in a goroutine for local development
	go func() {
		if err := e.Start(configs.AppConfig.ServerPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Shutting down the server: %v", err)
		}
	}()

	// Graceful Shutdown: Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully")
}
