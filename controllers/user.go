package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/ksoichiro/record/db"
	"github.com/ksoichiro/record/forms"
	"github.com/ksoichiro/record/models"
	"golang.org/x/crypto/bcrypt"
)

type UserController struct{}

func (u UserController) UserCreate(c *gin.Context) {
	var json forms.UserCreateForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := json.Name
	hash, err := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	password := string(hash)
	fmt.Printf("%s / %s\n", name, password)
	db := db.GetDB()
	tx := db.Begin()
	user := models.User{Name: name, Password: password, CreatedAt: time.Now()}
	tx.Create(&user)
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

func (u UserController) UserLogin(c *gin.Context) {
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid name or password"})
		return
	}

	signBytes, err := ioutil.ReadFile("./jwtRS256.key")
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
		"nbf": time.Now().UTC().Unix(),
		"exp": time.Now().UTC().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(signKey)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
