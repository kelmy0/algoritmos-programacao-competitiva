package utils

import (
	"crypto/ed25519"
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

func GenerateAccessToken(userId, username, email string, permissions []string, privateKey ed25519.PrivateKey, issuer string, isEmployee bool, expire_time time.Time) (tokenId string, tokenString string, err error) {
	tokenId, _, tokenString, err = generateToken(userId, username, email, issuer, "", permissions, privateKey, isEmployee, expire_time, false)
	return tokenId, tokenString, err
}

func GenerateRefreshToken(userId string, privateKey ed25519.PrivateKey, issuer string, expire_time time.Time, familyId string) (tokenId string, finalFamilyId string, tokenString string, err error) {
	return generateToken(userId, "", "", issuer, familyId, nil, privateKey, false, expire_time, true)
}

func GeneratePreAuthToken(userId, issuer string, privateKey ed25519.PrivateKey, expireTime time.Time) (tokenId string, tokenString string, err error) {
	tokenId, _, tokenString, err = generateToken(userId, "", "", issuer, "", nil, privateKey, false, expireTime, false)
	return tokenId, tokenString, err
}

func ValidateAccessToken(tokenString string, publicKey ed25519.PublicKey, expectedIssuer string) (*Claims, error) {
	return validateToken(tokenString, publicKey, expectedIssuer)
}

func ValidateRefreshToken(tokenString string, publicKey ed25519.PublicKey, expectedIssuer string) (*Claims, error) {
	return validateToken(tokenString, publicKey, expectedIssuer)
}

func generateToken(
	userId, username, email, issuer, familyId string,
	permissions []string,
	privateKey ed25519.PrivateKey,
	isEmployee bool,
	expire_time time.Time,
	isRefreshToken bool,
) (tokenId string, finalFamilyId string, tokenString string, err error) {

	if userId == "" || privateKey == nil || len(privateKey) != ed25519.PrivateKeySize || issuer == "" || expire_time.IsZero() {
		return "", "", "", errors.New("invalid token parameters: fields cannot be empty or zero")
	}

	tokenId, err = GenerateCustomId(32)
	if err != nil {
		return "", "", "", errors.New("error generating id token")
	}

	if isRefreshToken {
		if familyId == "" {
			finalFamilyId, err = GenerateCustomId(32)
			if err != nil {
				return "", "", "", errors.New("error generating family id")
			}
		} else {
			finalFamilyId = familyId
		}
	}

	claims := Claims{
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

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	tokenString, err = token.SignedString(privateKey)
	if err != nil {
		return "", "", "", err
	}

	return tokenId, finalFamilyId, tokenString, nil
}

func validateToken(tokenString string, publicKey ed25519.PublicKey, expectedIssuer string) (*Claims, error) {
	if publicKey == nil || len(publicKey) != ed25519.PublicKeySize {
		return nil, errors.New("invalid public key")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, errors.New("unexpected signature method")
		}
		return publicKey, nil
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
