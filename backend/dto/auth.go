package dto

import (
	"cloud.google.com/go/auth/credentials/idtoken"
)

type AuthRequest struct {
	Email    string `json:"email" binding:"required,min=5"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	Requires2FA  bool   `json:"requires_2fa"`
	PreAuthToken string `json:"pre_auth_token,omitempty"`
}

type Verify2FARequest struct {
	PreAuthToken string `json:"pre_auth_token" binding:"required"`
	Code         string `json:"code" binding:"required,len=6"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}

type GithubUserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Login string `json:"login"`
}

type GithubEmailResponse struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

type GoogleUserPayload struct {
	Subject       string
	Email         string
	Name          string
	EmailVerified bool
}

func NewGoogleUserPayload(payload *idtoken.Payload) (*GoogleUserPayload, string) {
	if payload == nil {
		return nil, CodeInternalError
	}

	email, ok := payload.Claims["email"].(string)
	if !ok || email == "" {
		return nil, CodeMissingGoogleEmail
	}

	emailVerified, ok := payload.Claims["email_verified"].(bool)
	if !ok || !emailVerified {
		return nil, CodeUnverifiedGoogleEmail
	}

	name, _ := payload.Claims["name"].(string)

	return &GoogleUserPayload{
		Subject:       payload.Subject,
		Email:         email,
		Name:          name,
		EmailVerified: emailVerified,
	}, ""
}
