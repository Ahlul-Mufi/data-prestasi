package utils

import (
	"errors"
	"os"
	"time"

	models "github.com/Ahlul-Mufi/data-prestasi/app/model/postgre"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	JwtKey     = []byte(getEnv("JWT_SECRET", "SUPER_SECRET_KEY"))
	RefreshKey = []byte(getEnv("REFRESH_SECRET", "SUPER_REFRESH_SECRET_KEY"))

	AccessTokenExpiry  = 15 * time.Minute
	RefreshTokenExpiry = 7 * 24 * time.Hour
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func GenerateToken(user models.User) (string, error) {
	now := time.Now()
	claims := models.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RoleID:   user.RoleID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

func ValidateToken(tokenString string) (*models.JWTClaims, error) {
	claims := &models.JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return JwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired, please use refresh token or login again")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GenerateRefreshToken(user models.User) (string, error) {
	now := time.Now()
	claims := models.RefreshClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(RefreshKey)
}

func ValidateRefreshToken(tokenString string) (*models.RefreshClaims, error) {
	claims := &models.RefreshClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return RefreshKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("refresh token has expired, please login again")
		}
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	return claims, nil
}

func GetTokenExpiry() time.Duration {
	return AccessTokenExpiry
}

func GetRefreshTokenExpiry() time.Duration {
	return RefreshTokenExpiry
}
