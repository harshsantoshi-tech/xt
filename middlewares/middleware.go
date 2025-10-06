package middlewares

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid Authorization header")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		// TODO: Validate token (add your real token verification logic here)
		if token != "your_valid_token" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		// Optionally, store user info in context for handler use
		c.Set("user", "demo_user") // Example

		return next(c)
	}
}
