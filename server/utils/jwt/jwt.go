package jwt

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID uint   `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}

func GenerateTokens(userID uint, login string) (accessToken string, refreshToken string, err error) {
	accessClaims := Claims{
		UserID: userID,
		Login:  login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * time.Minute)),
		},
	}
	accessToken, err = generateToken(accessClaims, []byte(os.Getenv("JWT_ACCESS_TOKEN")))
	if err != nil {
		return "", "", err
	}

	refreshClaims := Claims{
		UserID: userID,
		Login:  login,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	refreshToken, err = generateToken(refreshClaims, []byte(os.Getenv("JWT_REFRESH_TOKEN")))
	return
}

func generateToken(claims Claims, secret []byte) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ValidateAccessToken(tokenString string) (*Claims, error) {
	return validateTokenWithKey(tokenString, []byte(os.Getenv("JWT_ACCESS_TOKEN")))
}

func ValidateRefreshToken(tokenString string) (*Claims, error) {
	return validateTokenWithKey(tokenString, []byte(os.Getenv("JWT_REFRESH_TOKEN")))
}

func validateTokenWithKey(tokenString string, secretKey []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("неверный токен")
	}

	return claims, nil
}
