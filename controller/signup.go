package controller

import (
	"expense-tracker/constants"
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

	if !isValidEmail(req.Email) {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid email format",
			Code:    http.StatusBadRequest,
		})
	}

	err := handlers.SendOTPHandler(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseModel{
			Status:  "INTERNAL_SERVER_ERROR",
			Message: "Could not send OTP",
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  constants.SUCCESS,
		Message: "OTP sent to your email",
		Code:    http.StatusOK,
	})
}

func VerifySignupController(c echo.Context) error {
	var req models.VerifyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid input data",
			Code:    http.StatusBadRequest,
		})
	}

	if !isValidEmail(req.Email) {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid email format",
			Code:    http.StatusBadRequest,
		})
	}
	err := handlers.VerifySignupHandler(req.Email, req.OTP)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, models.ResponseModel{
			Status:  constants.UNAUTHORIZED,
			Message: err.Error(),
			Code:    http.StatusUnauthorized,
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  constants.SUCCESS,
		Message: "User registered successfully!",
		Code:    http.StatusOK,
	})
}
