package middlewares

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func UserAuthenticator() gin.HandlerFunc {
	return func(c *gin.Context) {
		verifyBytes, err := ioutil.ReadFile("./jwtRS256.key.pub.pkcs8")
		if err != nil {
			panic(err)
		}

		verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
		if err != nil {
			panic(err)
		}

		token, err := request.ParseFromRequest(c.Request, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
			_, err := token.Method.(*jwt.SigningMethodRSA)
			if !err {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return verifyKey, nil
		})
		claims, ok := token.Claims.(jwt.MapClaims)
		if err != nil || !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("user", claims["sub"])
		c.Next()
	}
}
