package middlewares

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/api/config"
)

// UserAuthenticator is a middleware to validate that the requests
// are authenticated with Authorization header.
func UserAuthenticator() gin.HandlerFunc {
	return func(c *gin.Context) {
		verifyBytes, err := ioutil.ReadFile(config.GetConfig().GetString("auth.keys.public"))
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
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}
		c.Set("user", claims["sub"])
		c.Next()
	}
}
