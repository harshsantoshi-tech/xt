package controller

import (
	"encoding/json"
	"expense-tracker/constants"
	"expense-tracker/handlers"
	"expense-tracker/models"
	"expense-tracker/services/redis"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/gommon/log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// Simple regex email validation
func isValidEmail(email string) bool {

	var emailRegex = regexp.MustCompile(`(?i)^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$`)
	email = strings.TrimSpace(email)
	return emailRegex.MatchString(email)

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

	if !isValidEmail(req.Email) {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid email format",
			Code:    http.StatusBadRequest,
		})
	}

	requestJson, _ := json.Marshal(req)
	log.Info("/login ", string(requestJson))

	accessToken, refreshToken, err := handlers.LoginHandler(req.Email, req.Password)
	if err != nil {
		log.Error("LoginController.LoginHandler ", err.Error())
		return c.JSON(http.StatusUnauthorized, models.ResponseModel{
			Status:  "UNAUTHORIZE",
			Message: "Invalid request",
			Code:    http.StatusUnauthorized,
		})
	}

	var response models.LoginResponse
	response.AccessToken = accessToken
	response.RefreshToken = refreshToken
	response.ResponseModel = models.ResponseModel{
		Status:  constants.SUCCESS,
		Message: "Logged in successfully",
		Code:    http.StatusOK,
	}
	return c.JSON(http.StatusOK, response)
}

func ForgotPasswordController(c echo.Context) error {

	var req models.ForgotPasswordRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid request",
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

	if !isValidEmail(req.Email) {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid email format",
			Code:    http.StatusBadRequest,
		})
	}

	pending, err := redis.GetPendingUser(req.Email)
	if err != nil || pending.OTP != req.OTP {
		log.Error("VerifyResetOTPController.VerifyResetOTPHandler ", " OTP Mismatch ")
		return c.JSON(http.StatusUnauthorized, models.ResponseModel{
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
	var req models.ResetPasswordRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid request",
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
	if req.Password != req.ConfirmPassword {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Confirm password does not match",
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

func ResendOTPController(c echo.Context) error {
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

	if !isValidEmail(req.Email) {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  "BAD_REQUEST",
			Message: "Invalid email format",
			Code:    http.StatusBadRequest,
		})
	}

	err := handlers.ResendOTPHandler(req.Email)
	if err != nil {
		log.Error("ResendOTPController.ResendOTPHandler ", err.Error())
		return c.JSON(http.StatusInternalServerError, models.ResponseModel{
			Status:  constants.FAILURE,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  constants.SUCCESS,
		Message: "A new OTP has been sent to your email",
		Code:    http.StatusOK,
	})
}

func LogoutController(c echo.Context) error {

	authHeader := c.Request().Header.Get("Authorization")
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  constants.BAD_REQUEST,
			Message: "Invalid or missing token",
			Code:    http.StatusBadRequest,
		})
	}
	tokenString := authHeader[7:]

	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.ResponseModel{
			Status:  constants.BAD_REQUEST,
			Message: "Invalid token",
			Code:    http.StatusBadRequest,
		})
	}

	var expiration time.Duration
	if claims, ok := token.Claims.(jwt.MapClaims); ok && claims["exp"] != nil {
		exp := int64(claims["exp"].(float64))
		expiryTime := time.Unix(exp, 0)
		expiration = time.Until(expiryTime)
	} else {
		expiration = constants.JWT_EXPIRY_TIME
	}

	err = redis.SetRedis(fmt.Sprintf(redis.USER_LOGOUT, tokenString), "true", expiration)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseModel{
			Status:  constants.FAILURE,
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, models.ResponseModel{
		Status:  constants.SUCCESS,
		Message: "Logged out successfully",
		Code:    http.StatusOK,
	})
}