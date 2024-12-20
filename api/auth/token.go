package auth

import (
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	Admin  bool   `json:"admin"`
	jwt.RegisteredClaims
}

func NewToken(key []byte, userID string, admin bool) (string, error) {
	claims := CustomClaims{
		userID,
		admin,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(14 * 24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(key)
}

func NewConfirmToken() string {
	return gofakeit.DigitN(6)
}

func NewResetToken() string {
	return gofakeit.DigitN(6)
}
