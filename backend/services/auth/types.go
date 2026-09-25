package auth

import "github.com/kelmy0/algoritmos-programacao-competitiva/backend/dto"

type AuthResult struct {
	LoginResponse dto.LoginResponse
	RefreshToken  string
}

type RefreshTokenResult struct {
	AccessToken  string
	RefreshToken string
}

type SignUpResult struct {
	SignUpResponse dto.SignUpResponse
	RefreshToken   string
}
