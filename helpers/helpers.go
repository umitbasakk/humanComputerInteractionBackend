package helpers

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateJWTToken(username string) (string, error) {
	jwtHash := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
		},
	)

	tokenString, err := jwtHash.SignedString([]byte("WvHtwWc2ctdCFzdQUv3ZHmuPB12fVxtb"))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", tokenString), nil
}

func ValidEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}
