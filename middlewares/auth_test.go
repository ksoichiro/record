package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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
	req.Header.Add("Authorization", "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJodHRwczovL2FwaS5leGFtcGxlLmNvbSIsImV4cCI6MTU3NDM0ODMzMSwiaXNzIjoiaHR0cHM6Ly9pZHAuZXhhbXBsZS5jb20iLCJuYmYiOjE1NzQyNjE5MzEsInN1YiI6MX0.WBTXuWwyWHJzpJsIR5dqpmbDVuPa90sLE7xGhe_uMw8TXSrdUwxCUhmPU3xgPq2_o0W0MjaeofKhsUelD-fcGa23KqwC-_sjN8P7IaNwnNaaRMwdqx-b6EZh6s4p-lRPRWzAu3ldDQL5JSZk1s5uiZZZBUp7cXBerb7BB_DwVQYBy7pkzZ_q-IBYV7ik44V8oEEBLa1I7WxhGL60PkM_-P6HbV8HZwAozwvWwXoKpiqTfKcFO37cwEQO_fF0gtjK-4hH3uLab5Dcw7mpSFgefiI3wH6dk7TcNmV69a_hKhpffg04AzGXvaB-NUBHAiEvTCubxWhytGFdYPwMQZfTvko5zHY2KAjA-xkcwWp7bC77vGqh1kM2UjnOl6cwGTWfHH9vrA7WfgKQeW4P0-_NpVD5v4eGx91a4qYzqZw3e1sE6jQKzpbaMjqEavuncfERFWlx8lni4zDyXzHCfBC41wpghAG3706QJy_3s2ePTZL2dVYH1cXiEy-ritPXgBnGJ7nX0Sw10zeGsuPZ60ujqlcwxtwXzx1yYYXLNi3MWx_e-GD5hOfG0NDTxQUXXNsBaLMMo624D0C_AyKQ67hbjad4QhDetgnG0QoYT5bX9HkpNK0U2NRD6l42oi2JrLQsEdO8efRuxBf_69YswZ9WByQ85eh4GzOeL7EFSvb_cU8")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}
