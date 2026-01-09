package models

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Helper functions for common errors
func InternalError(msg string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: msg}
}

func BadRequest(msg string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: msg}
}

func Unauthorized(msg string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: msg}
}
