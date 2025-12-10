package utils

import (
	"errors"
	"strings"
	"time"

	models "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("SUPER_SECRET_KEY")

func GenerateToken(user models.User) (string, error) {
	claims := models.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
        RoleID:   user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey) 
}

func ValidateToken(tokenString string) (*models.JWTClaims, error) {
	claims := &models.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		return nil, err
	}

    if !token.Valid {
        return nil, errors.New("invalid token")
    }

    if claims, ok := token.Claims.(*models.JWTClaims); ok && token.Valid {
        return claims, nil
    }

	return nil, errors.New("invalid token or key") 
}

func ExtractToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("missing authorization header")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}