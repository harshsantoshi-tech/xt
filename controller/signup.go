package controller

import (
	"expense-tracker/handlers"
	"expense-tracker/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

func SignupController(c echo.Context) error {
	var req models.SignupRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid input data",
			Code:    http.StatusBadRequest,
		})
	}

	// Call handler to generate and send OTP
	// This usually saves the user info + OTP temporarily (in Redis)
	err := handlers.SendOTPHandler(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseModel{
			Status:  "INTERNAL_SERVER_ERROR",
			Message: "Could not send OTP",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  "SUCCESS",
		Message: "OTP sent to your email",
		Code:    http.StatusOK,
	})
}

func VerifySignupController(c echo.Context) error {
	var req models.VerifyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request"})
	}

	err := handlers.VerifySignupHandler(req.Email, req.OTP)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ResponseModel{
			Status:  "UNAUTHORIZED",
			Message: err.Error(),
			Code:    http.StatusUnauthorized,
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "User registered successfully!",
	})
}
