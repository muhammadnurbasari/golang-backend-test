package middleware

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	jwtSigningMethod = jwt.SigningMethodHS256
	jwtSignatureKey  = []byte(os.Getenv("APP_SECRET_KEY_JWT"))
)

func BasicAuth(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")

	if !strings.Contains(authHeader, "Basic") {
		result := gin.H{
			"status":  http.StatusForbidden,
			"message": "Invalid Token",
			"href":    c.Request.RequestURI,
		}

		c.JSON(http.StatusForbidden, result)
		c.Abort()
		return
	}

	clientId := os.Getenv("CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")

	tokenString := strings.Replace(authHeader, "Basic ", "", -1)
	myToken := clientId + ":" + clientSecret
	myBasicToken := base64.StdEncoding.EncodeToString([]byte(myToken))

	if tokenString != myBasicToken {
		result := gin.H{
			"status":  http.StatusUnauthorized,
			"message": "Unauthorized user",
			"href":    c.Request.RequestURI,
		}
		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

}

func JwtAuth(c *gin.Context) {
	authHeader := c.Request.Header.Get("authorization")
	if !strings.Contains(authHeader, "Bearer") {
		result := gin.H{
			"status":  http.StatusForbidden,
			"message": "Invalid Token",
			"href":    c.Request.RequestURI,
		}

		c.JSON(http.StatusForbidden, result)
		c.Abort()
		return
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", -1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate algoritma its used
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else if method != jwtSigningMethod {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSignatureKey, nil
	})

	if err != nil {
		result := gin.H{
			"status":  http.StatusUnauthorized,
			"message": err.Error(),
		}

		c.JSON(http.StatusUnauthorized, result)
		c.Abort()
		return
	}

	_, okToken := token.Claims.(jwt.MapClaims)

	if !okToken {
		result := gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Error check payload token",
		}

		c.JSON(http.StatusInternalServerError, result)
		c.Abort()
		return
	}

}
