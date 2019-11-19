package main

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

type user struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type userCreateForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userLoginForm userCreateForm

type task struct {
	ID     int
	UserID int
	Title  string
}

func connect() *gorm.DB {
	db, err := gorm.Open("mysql", "app:app@tcp(127.0.0.1:3306)/test?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		db := connect()
		defer db.Close()
		task := task{}
		db.First(&task)
		fmt.Printf("%d: UserId: %d, Title: %s\n", task.ID, task.UserID, task.Title)

		c.String(http.StatusOK, "pong")
	})

	userGroup := r.Group("/user")
	userGroup.POST("/create", userCreate)
	userGroup.POST("/login", userLogin)

	taskGroup := r.Group("/task")
	taskGroup.Use(userAuthenticator())
	taskGroup.POST("/create", taskCreate)

	adminGroup := r.Group("/admin")
	adminGroup.GET("/user", adminUser)

	return r
}

func userCreate(c *gin.Context) {
	var json userCreateForm
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
	db := connect()
	defer db.Close()
	tx := db.Begin()
	user := user{Name: name, Password: password, CreatedAt: time.Now()}
	tx.Create(&user)
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

func userLogin(c *gin.Context) {
	var json userLoginForm
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := json.Name

	db := connect()
	defer db.Close()

	user := user{}
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

func userAuthenticator() gin.HandlerFunc {
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
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func taskCreate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "created"})
}

func adminUser(c *gin.Context) {
	db := connect()
	defer db.Close()
	users := []user{}
	db.Find(&users)
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
