package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	Permissions []string `json:"permissions"`
	IsEmployee  bool     `json:"isEmployee"`
	FamilyId    string   `json:"familyId,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(userId, username, email string, permissions []string, secretKey, issuer string, isEmployee bool, expire_time time.Time, familyId string) (tokenId string, finalFamilyId string, tokenString string, err error) {
	if userId == "" || secretKey == "" || issuer == "" || expire_time.IsZero() {
		return "", "", "", errors.New("invalid token parameters: fields cannot be empty or zero")
	}

	tokenId, err = GenerateCustomId(32)
	if err != nil {
		return "", "", "", errors.New("Error generating id token")
	}

	if familyId == "" {
		finalFamilyId, err = GenerateCustomId(32)
		if err != nil {
			return "", "", "", errors.New("error generating family id")
		}
	} else {
		finalFamilyId = familyId
	}

	claimsRefresh := Claims{
		Username:    username,
		Email:       email,
		Permissions: permissions,
		IsEmployee:  isEmployee,
		FamilyId:    finalFamilyId,
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
	tokenString, err = token.SignedString(secretByte)
	if err != nil {
		return "", "", "", err
	}

	return tokenId, finalFamilyId, tokenString, nil
}

func ValidateToken(tokenString, secretKey, expectedIssuer string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signature method")
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
		return nil, errors.New("invalid token issuer")
	}

	return claims, nil
}
