package helpers

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(getJWTSecret())

func getJWTSecret() string {
	if secret := os.Getenv("JWT_SECRET"); secret != "" {
		return secret
	}
	return "YOUR_ULTRA_SECRET_KEY_CHANGE_ME"
}

type Claims struct {
	UserID   uint   `json:"user_id"`
	RoleName string `json:"role_name"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, roleName string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID:   userID,
		RoleName: roleName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}
	return claims, nil
}
