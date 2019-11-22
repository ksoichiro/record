package middlewares

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/models"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestUserAuthenticator(t *testing.T) {
	router := gin.Default()
	gin.SetMode(gin.TestMode)
	db.InitForTest()
	db := db.GetDB()
	db.AutoMigrate(&models.User{}, &models.Task{})
	db.Create(&models.User{ID: 100, Name: "foo", Password: "$2a$10$FgKFrUubZOpRwPT9D5p9XuOjCYhPv7eCQwzdQKFJWTQsC9tXAuMG2" /* test */, CreatedAt: time.Now()})
	resourcePath = ".."
	router.Use(UserAuthenticator())
	router.GET("/", func(c *gin.Context) {})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	signBytes, err := ioutil.ReadFile("../jwtRS256.key")
	if err != nil {
		panic(err)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "https://idp.example.com",
		"aud": "https://api.example.com",
		"sub": 100,
		"nbf": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", tokenString)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
