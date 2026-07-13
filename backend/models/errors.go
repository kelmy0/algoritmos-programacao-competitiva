package models

import (
	"errors"
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"error"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrPasswordsDontMatch         = &AppError{StatusCode: http.StatusBadRequest, Code: "USER_PASSWORDS_DONT_MATCH", Message: "Passwords don't match."}
	ErrIncorrectPassword          = &AppError{StatusCode: http.StatusBadRequest, Code: "AUTH_INCORRECT_PASSWORD", Message: "Incorrect password."}
	ErrUserNotFound               = &AppError{StatusCode: http.StatusNotFound, Code: "USER_NOT_FOUND", Message: "User not found."}
	ErrInvalidOrExpiredRefresh    = &AppError{StatusCode: http.StatusUnauthorized, Code: "AUTH_INVALID_REFRESH_TOKEN", Message: "Invalid or expired refresh token."}
	ErrTokenMetadataMisMatch      = &AppError{StatusCode: http.StatusForbidden, Code: "AUTH_SECURITY_VIOLATION", Message: "Token metadata mismatch: Security violation."}
	ErrPasswordChangeFailed       = &AppError{StatusCode: http.StatusInternalServerError, Code: "USER_PASSWORD_CHANGE_FAILED", Message: "Failed to change password."}
	ErrPasswordNotSet             = &AppError{StatusCode: http.StatusBadRequest, Code: "USER_PASSWORD_NOT_SET", Message: "Password not set before."}
	ErrAlgorithmNotFound          = &AppError{StatusCode: http.StatusNotFound, Code: "ALGORITHM_NOT_FOUND", Message: "Algorithm not found."}
	ErrFailQueryingAlgorithm      = &AppError{StatusCode: http.StatusInternalServerError, Code: "ALGORITHM_QUERY_FAILED", Message: "Failed to query algorithm."}
	ErrInvalidNameCategoryContent = &AppError{StatusCode: http.StatusBadRequest, Code: "ALGORITHM_INVALID_NAME_CATEGORY_CONTENT", Message: "Invalid name, content or category."}
	ErrFailGeneratePublicId       = &AppError{StatusCode: http.StatusInternalServerError, Code: "ALGORITHM_GENERATE_PUBLIC_ID_FAILED", Message: "Failed to generate public ID."}
	ErrFailQueryUser              = &AppError{StatusCode: http.StatusInternalServerError, Code: "AUTH_QUERY_USER_FAILED", Message: "Failed to query the user in the database."}
	ErrRegisterSocialUser         = &AppError{StatusCode: http.StatusInternalServerError, Code: "USER_REGISTER_SOCIAL_FAILED", Message: "Failed to register social user."}
	ErrLinkGoogleAccount          = &AppError{StatusCode: http.StatusInternalServerError, Code: "LINK_GOOGLE_ACCOUNT_FAILED", Message: "Failed to link google account."}
	ErrReloadUser                 = &AppError{StatusCode: http.StatusInternalServerError, Code: "USER_RELOAD_FAILED", Message: "Failed to reload user after link."}
	ErrUserNotEnabled             = &AppError{StatusCode: http.StatusForbidden, Code: "USER_NOT_ENABLED", Message: "This account is not enabled."}
	ErrUnexpectedLogin            = &AppError{StatusCode: http.StatusInternalServerError, Code: "AUTH_UNEXPECTED_ERROR", Message: "Unexpected login error."}
	ErrGeneratingToken            = &AppError{StatusCode: http.StatusInternalServerError, Code: "AUTH_GENERATE_TOKEN_FAILED", Message: "Error generating Token."}
	ErrSessionExpired             = &AppError{StatusCode: http.StatusUnauthorized, Code: "SESSION_EXPIRED", Message: "Session expired."}
	ErrInvalidEmailOrPassword     = &AppError{StatusCode: http.StatusUnauthorized, Code: "AUTH_INVALID_EMAIL_PASSWORD", Message: "Invalid email or password."}
	ErrSessionData                = &AppError{StatusCode: http.StatusBadRequest, Code: "INVALID_SESSION_DATA", Message: "Invalid session data."}
	Err2FANotInitiated            = &AppError{StatusCode: http.StatusPreconditionFailed, Code: "2FA_NOT_INITIATED", Message: "2FA setup has not been initiated for this user."}
	Err2FAInvalid                 = &AppError{StatusCode: http.StatusUnauthorized, Code: "2FA_INVALID_CODE", Message: "2FA code is invalid or expired."}
	ErrUserAlreadyExists          = &AppError{StatusCode: http.StatusConflict, Code: "USER_ALREADY_EXISTS", Message: "User already exists."}
	ErrInvalidRegistrationFields  = &AppError{StatusCode: http.StatusBadRequest, Code: "REGISTRATION_INVALID_FIELDS", Message: "Invalid fields."}
	ErrInvalidEmailFormat         = &AppError{StatusCode: http.StatusBadRequest, Code: "REGISTRATION_INVALID_EMAIL", Message: "Invalid email format."}
	ErrUserRegistrationFailed     = &AppError{StatusCode: http.StatusInternalServerError, Code: "REGISTRATION_UNEXPECTED_ERROR", Message: "Unexpected registration error."}
)

var ErrPasswordChangeButNotLogout = errors.New("password changed, but failed to log out other devices.")
var ErrAccountCreatedButTokenFailed = errors.New("account created successfully, but auto-login failed")
