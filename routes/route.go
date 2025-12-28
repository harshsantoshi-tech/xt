package routes

import (
	"expense-tracker/controller"
	"expense-tracker/handlers"
	"expense-tracker/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// InitRoutes sets up all the API routes for the application.
func InitRoutes(e *echo.Echo) {
	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Public route
	e.GET("/", handlers.GetRoot)

	// Authenticated routes group
	authGroup := e.Group("/api")
	authGroup.Use(middlewares.AuthMiddleware) // Apply your custom auth middleware

	// Protected routes
	e.GET("/home", controller.HomePageController)
	e.GET("/ws", controller.WsController)
	e.GET("/text/change", handlers.TextChangeHandler)
}
