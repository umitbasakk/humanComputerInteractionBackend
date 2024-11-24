package helpers

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"net/mail"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomJWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func CreateJWTToken(username string) (string, error) {

	userClaim := &CustomJWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 100)),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaim)

	signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", signedAccessToken), nil
}

func ParseJWT(token string) (*CustomJWTClaims, error) {
	parsedJwtAccessToken, err := jwt.ParseWithClaims(token, &CustomJWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	} else if claims, ok := parsedJwtAccessToken.Claims.(*CustomJWTClaims); ok {
		return claims, nil
	} else {
		return nil, errors.New("unkown claims type")
	}
}

func IsClaimExpired(claims *CustomJWTClaims) bool {
	currentTime := jwt.NewNumericDate(time.Now())
	return claims.ExpiresAt.Time.Before(currentTime.Time)
}

func ValidEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func GetVerifyCode() int {
	return rand.IntN(9999-1000) + 1000
}
