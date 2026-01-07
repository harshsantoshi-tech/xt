package controller

import (
	"encoding/json"
	"expense-tracker/constants"
	"expense-tracker/handlers"
	"expense-tracker/models"
	"expense-tracker/services/redis"
	"fmt"
	"github.com/labstack/gommon/log"
	"net/http"
	"regexp"
	"time"

	"github.com/labstack/echo/v4"
)

// Simple regex email validation
func isValidEmail(email string) bool {
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

func LoginController(c echo.Context) error {
	var req models.LoginRequest

	if err := c.Bind(&req); err != nil {
		log.Error("LoginController.Bind ", err.Error())
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	requestJson, _ := json.Marshal(req)
	log.Info("/login ", string(requestJson))

	token, err := handlers.LoginHandler(req.Email, req.Password)
	if err != nil {
		log.Error("LoginController.LoginHandler ", err.Error())
		return c.JSON(http.StatusUnauthorized, models.ResponseModel{
			Status:  "UNAUTHORIZE",
			Message: "Invalid request",
			Code:    http.StatusUnauthorized,
		})
	}

	var response models.LoginResponse
	response.Token = token
	response.ResponseModel = models.ResponseModel{
		Status:  constants.SUCCESS,
		Message: "Logged in successfully",
		Code:    http.StatusOK,
	}
	return c.JSON(http.StatusOK, response)
}

func ForgotPasswordController(c echo.Context) error {
	var req struct {
		Email string `json:"email"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	err := handlers.ForgotPasswordHandler(req.Email)
	if err != nil {
		log.Error("ForgotPasswordController.ForgotPasswordHandler ", err.Error())
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  constants.SUCCESS,
		Message: "OTP sent if account exists ",
		Code:    http.StatusOK,
	})
}

func VerifyResetOTPController(c echo.Context) error {

	var req models.VerifyRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	pending, err := redis.GetPendingUser(req.Email)
	if err != nil || pending.OTP != req.OTP {
		log.Error("VerifyResetOTPController.VerifyResetOTPHandler ", " OTP Mismatch ")
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "UNAUTHORIZED",
			Message: "Invalid OTP",
			Code:    http.StatusUnauthorized,
		})
	}

	redis.SetRedis(fmt.Sprintf(redis.RESET_ALLOWED, req.Email), "true", 10*time.Minute)
	redis.ClearSignupData(req.Email)

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  constants.SUCCESS,
		Message: "OTP verified. You may now reset password.",
		Code:    http.StatusOK,
	})
}

func ResetPasswordController(c echo.Context) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid request",
			Code:    http.StatusBadRequest,
		})
	}

	err := handlers.ResetPasswordHandler(req.Email, req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseModel{
			Status:  constants.FAILURE,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  constants.SUCCESS,
		Message: "Password updated successfully",
		Code:    http.StatusOK,
	})
}
