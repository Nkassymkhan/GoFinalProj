package main

import (
"net/http"
"os"
"strings"
"time"

"github.com/dgrijalva/jwt-go"
"github.com/gin-gonic/gin"
"github.com/jinzhu/gorm"
_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Product struct {
ID    uint   `json:"id"`
Name  string `json:"name"`
Price uint   `json:"price"`
}

var db *gorm.DB

func init() {
var err error
db, err = gorm.Open("sqlite3", "test.db")
if err != nil {
panic("failed to connect database")
}

db.AutoMigrate(&Product{})
}

func login(c *gin.Context) {
	var user struct {
	Username string `json:"username"`
	Password string `json:"password"`
	}

	if err := c.BindJSON(&user); err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
	return
	}

	if user.Username == "admin" && user.Password == "password" {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"username": user.Username,
	"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	return
}

c.JSON(http.StatusOK, gin.H{"token": tokenString})
} else {
c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
}
}

func authMiddleware() gin.HandlerFunc {
return func(c *gin.Context) {
authHeader := c.GetHeader("Authorization")
if authHeader == "" {
c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
c.Abort()
return
}

tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	return []byte(os.Getenv("SECRET_KEY")), nil
	})

	if err != nil || !token.Valid {
	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
	c.Abort()
	return
	}

	c.Next()
	}
}

func createProduct(c *gin.Context) {
	var newProduct Product
	if err := c.BindJSON(&newProduct); err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Input is not correct"})
	return
	}

	db.Create(&newProduct)
	c.JSON(http.StatusOK, newProduct)
}


