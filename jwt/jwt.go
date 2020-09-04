package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	appName          = os.Getenv("APP_NAME")
	jwtSigningMethod = jwt.SigningMethodHS256
	jwtSigningKey    = []byte(os.Getenv("APP_SECRET_KEY_JWT"))
)

// MyCustomClaims - object to custom claims
type MyCustomClaims struct {
	UserId int `json:"user_id"`
	jwt.StandardClaims
}

// GenerateTokenJwt - generate token JWT
func GenerateTokenJwt(userId int, timeExpDuration time.Duration) (string, error) {
	// Create the Claims
	claims := MyCustomClaims{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(timeExpDuration).Unix(),
			Issuer:    appName,
		},
	}

	token := jwt.NewWithClaims(jwtSigningMethod, claims)

	signingTokenString, err := token.SignedString(jwtSigningKey)

	if err != nil {
		return "", errors.New("Generate JWT err = " + err.Error())
	}

	return signingTokenString, nil
}
