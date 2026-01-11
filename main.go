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
)

// A global Echo instance to be reused across requests.
var e *echo.Echo

//// Handler is the Vercel-compatible entry point for your serverless function.
//// It simply serves the incoming request using the global Echo instance.
//func Handler(w http.ResponseWriter, r *http.Request) {
//	e.ServeHTTP(w, r)
//}


func init() {
	e = echo.New()
	configs.InitializeConfigs()
	routes.InitRoutes(e)
}

func main() {
	// Start the server in a goroutine for local development
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default for local testing
	}
	// Start server
	go func() {
		if err := e.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
