package application

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type MyCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

// 🔐 Obtener secret seguro
func getSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}
	return secret
}

// ✅ GENERAR TOKEN REAL
func GenerateTokenWithEmail(email string) string {
	claims := MyCustomClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(getSecret()))
	if err != nil {
		log.Println(err)
	}

	return signedToken
}

// ✅ VALIDAR TOKEN
func ValidateToken(authorizationToken string) (string, error) {
	receivedToken := strings.Replace(authorizationToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(receivedToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(getSecret()), nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims.Email, nil
	}

	return "", fmt.Errorf("invalid token")
}
