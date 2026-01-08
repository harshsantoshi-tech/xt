package middlewares

import (
	"expense-tracker/constants"
	"expense-tracker/models"
	"expense-tracker/services"
	"expense-tracker/services/redis"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header missing")
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header format")
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			return echo.NewHTTPError(http.StatusInternalServerError, "JWT secret not configured")
		}

		userLoggedOut, err := redis.GetRedis(fmt.Sprintf(redis.USER_LOGOUT, tokenString))

		if err != nil || userLoggedOut == "true" {
			return c.JSON(http.StatusUnauthorized, models.ResponseModel{
				Status:  constants.UNAUTHORIZED,
				Message: "Token has been invalidated. Please login again.",
				Code:    http.StatusUnauthorized,
			})
		}

		userId, err := services.GetUserIDFromToken(tokenString, secret)
		if err != nil {
			log.Error("Error getting user id from token ", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "Error getting user id from token")
		}

		c.Set("user_id", userId)

		return next(c)
	}
}
