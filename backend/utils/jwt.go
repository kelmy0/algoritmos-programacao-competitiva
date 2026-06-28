package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleId   int    `json:"role_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userId, username, email string, roleId int, secretKey, issuer string, expire_time time.Time) (string, error) {
	tokenId, err := GenerateCustomId(32)

	if err != nil {
		return "", errors.New("Error generating id token")
	}

	claimsRefresh := Claims{
		Username: username,
		Email:    email,
		RoleId:   roleId,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenId,
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(expire_time),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}
	secretByte := []byte(secretKey)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	return token.SignedString(secretByte)
}
