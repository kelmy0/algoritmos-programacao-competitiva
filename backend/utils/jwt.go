package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	RoleId int    `json:"role_id"`
	jwt.RegisteredClaims
}

func GenerateAcessToken(userId string, email string, roleId int, secretKey string, issuer string) (string, error) {
	claimsAccess := Claims{
		UserId: userId,
		Email:  email,
		RoleId: roleId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer, // por um NAME no .env depois
		},
	}

	secretByte := []byte(secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsAccess)
	return token.SignedString(secretByte)
}

func GenerateRefreshToken(userId string, email string, roleId int, secretKey string, issuer string) (string, error) {
	claimsRefresh := Claims{
		UserId: userId,
		Email:  email,
		RoleId: roleId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 0, 7)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer, // por um NAME no .env depois
		},
	}

	secretByte := []byte(secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	return token.SignedString(secretByte)
}
