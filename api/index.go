package handler

import (
	"expense-tracker/routes"
	"net/http"

	"github.com/labstack/echo/v4"
)

// A global Echo instance to be reused across requests.
var e *echo.Echo

func init() {
	e = echo.New()
	//e.Logger.SetLevel(echo.IN)
	routes.InitRoutes(e)
}

// Handler is the Vercel-compatible entry point for your serverless function.
func Handler(w http.ResponseWriter, r *http.Request) {
	e.ServeHTTP(w, r)
}
