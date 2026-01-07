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
	// 1. Global Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS()) // Added CORS for frontend connectivity

	// 2. Public Routes (No Auth Required)
	e.GET("/", handlers.GetRoot)
	e.POST("/login", controller.LoginController)

	// New Signup Flow Routes
	e.POST("/signup-request", controller.SignupController)
	e.POST("/signup-verify", controller.VerifySignupController)

	//Forgot password
	e.POST("/forgot-password", controller.ForgotPasswordController)
	e.POST("/verify-otp", controller.VerifyResetOTPController)
	e.POST("/reset-password", controller.ResetPasswordController)

	// 3. Protected Routes Group
	api := e.Group("/api")
	api.Use(middlewares.AuthMiddleware) // Tokens required for everything below
	{
		api.GET("/home", controller.HomePageController)
		api.GET("/text/change", handlers.TextChangeHandler)
		api.GET("/ws", controller.WsController)
	}
}
