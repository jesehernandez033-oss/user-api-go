package application

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"strings"
	"time"
)

type MyCustomClaims struct {
	Email    string `json:"email"`
	Birthday int64  `json:"birthday"`
	jwt.RegisteredClaims
}

var mySecret = "my-secret"

func GenerateToken() string {
	claims := MyCustomClaims{
		"hello@friendsofgo.tech",
		time.Date(2019, 01, 01, 0, 0, 0, 0, time.UTC).Unix(),
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "test",
			Subject:   "somebody",
			ID:        "1",
			Audience:  []string{"somebody_else"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(mySecret))
	if err != nil {
		log.Println(err)
	}
	return signedToken
}

func ValidateToken(authorizationToken string) (string, error) {
	//r.Header.Get("Authorization")
	receivedToken := strings.Replace(authorizationToken, "Bearer ", "", -1)

	token, err := jwt.ParseWithClaims(receivedToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(mySecret), nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return fmt.Sprintf("%s,%d", claims.Email, claims.Birthday), nil
	}

	return "", err
}
