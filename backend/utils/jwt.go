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
	if userId == "" || username == "" || email == "" || secretKey == "" || issuer == "" || expire_time.IsZero() || roleId <= 0 {
		return "", errors.New("invalid token parameters: fields cannot be empty or zero")
	}

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

func ValidadeToken(tokenString, secretKey, expectedIssuer string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unespect signature method.")
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	if claims.Issuer != expectedIssuer {
		return nil, errors.New("invalid token issuer.")
	}

	return claims, nil
}
