package controllers

import (
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/config"
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/forms"
	"github.com/ksoichiro/record/models"
	"golang.org/x/crypto/bcrypt"
)

// UserController handles requests about users.
type UserController struct {
}

// Create creates a new user.
func (u UserController) Create(c *gin.Context) {
	var json forms.UserCreateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := models.NewUser(&json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

// Login authenticates the user and returns a authentication token.
func (u UserController) Login(c *gin.Context) {
	var json forms.UserLoginForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := json.Name

	db := db.GetDB()

	user := models.User{}
	db.Where("name = ?", name).First(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(json.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid name or password"})
		return
	}

	signBytes, err := ioutil.ReadFile(config.GetConfig().GetString("auth.keys.private"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		panic(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss": "https://idp.example.com",
		"aud": "https://api.example.com",
		"sub": user.ID,
		"nbf": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(signKey)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
