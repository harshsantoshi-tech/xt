package controller

import (
	"expense-tracker/handlers"
	"expense-tracker/models"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
)

// func SignupController(c echo.Context) error {
// 	req := new(models.LoginRequest)
// 	if err := c.Bind(req); err != nil {
// 		return c.JSON(http.StatusBadRequest, echo.Map{
// 			"error": "Invalid request format",
// 		})
// 	}

// 	// Sanitize input
// 	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

// 	// Basic validation
// 	if req.Email == "" || req.Password == "" {
// 		return c.JSON(http.StatusBadRequest, echo.Map{
// 			"error": "Email and password are required",
// 		})
// 	}
// 	if !isValidEmail(req.Email) {
// 		return c.JSON(http.StatusBadRequest, echo.Map{
// 			"error": "Invalid email format",
// 		})
// 	}
// 	if len(req.Password) < 6 {
// 		return c.JSON(http.StatusBadRequest, echo.Map{
// 			"error": "Password must be at least 6 characters",
// 		})
// 	}

// 	// Check for existing user
// 	_, err := services.GetUserFromEmail(req.Email)
// 	if err == nil {
// 		// User already exists
// 		return c.JSON(http.StatusConflict, echo.Map{
// 			"error": "Email already registered",
// 		})
// 	}

// 	// Hash the password
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return c.JSON(http.StatusInternalServerError, echo.Map{
// 			"error": "Failed to hash password",
// 		})
// 	}

// 	// Create user
// 	user := &models.User{
// 		Email:    req.Email,
// 		Password: string(hashedPassword),
// 	}
// 	if err := services.CreateUser(user); err != nil {
// 		return c.JSON(http.StatusInternalServerError, echo.Map{
// 			"error": "Failed to create user",
// 		})
// 	}

// 	// Optional: generate token here

// 	return c.JSON(http.StatusCreated, echo.Map{
// 		"message": "User registered successfully",
// 		"user": echo.Map{
// 			"id":    user.Id,
// 			"email": user.Email,
// 		},
// 	})
// }

// Simple regex email validation
func isValidEmail(email string) bool {
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}


func LoginController(c echo.Context) error {
	var req models.LoginRequest
	
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status: "BAD_REQUEST",
			Message: "",
			Code :http.StatusBadRequest,
		})
	}

	token, err := handlers.LoginHandler(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ResponseModel{
			Status: "UNAUTHORIZE",
			Message: "",
			Code :http.StatusUnauthorized,
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}