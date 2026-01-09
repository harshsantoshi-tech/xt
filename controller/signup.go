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
	if req.Password == "" || req.ConfirmPassword != req.Password {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid request",
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

	if appErr := handlers.VerifySignupHandler(req.Email, req.OTP); appErr != nil {
		return c.JSON(appErr.Code, models.ResponseModel{
			Status:  constants.FAILURE,
			Message: appErr.Message,
			Code:    appErr.Code,
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  constants.SUCCESS,
		Message: "User registered successfully!",
		Code:    http.StatusOK,
	})
}
