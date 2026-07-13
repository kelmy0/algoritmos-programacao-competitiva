package models

import (
	"errors"
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"error"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrPasswordsDontMatch      = &AppError{StatusCode: http.StatusBadRequest, Message: "Passwords don't match"}
	ErrIncorrectPassword       = &AppError{StatusCode: http.StatusBadRequest, Message: "Incorrect password"}
	ErrUserNotFound            = &AppError{StatusCode: http.StatusNotFound, Message: "User not found"}
	ErrInvalidOrExpiredRefresh = &AppError{StatusCode: http.StatusUnauthorized, Message: "Invalid or expired refresh token"}
	ErrTokenMetadataMisMatch   = &AppError{StatusCode: http.StatusForbidden, Message: "Token metadata mismatch: Security violation"}
	ErrFailChangingPassword    = &AppError{StatusCode: http.StatusInternalServerError, Message: "Fail changing password"}
	ErrPasswordNotSet          = &AppError{StatusCode: http.StatusBadRequest, Message: "Password not set before"}
)

var ErrPasswordChangeButNotLogout = errors.New("password changed, but failed to log out other devices")
