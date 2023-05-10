package authorization

import (
	"github.com/Nkassymkhan/GoFinalProj.git/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	db *gorm.DB
}

func New(db *gorm.DB) handler {
	return handler{db}
}

func (h * handler)Register(c *gin.Context) {
	var user models.User
	var db gorm.DB
	c.BindJSON(&user)
  
	if user.Username == "" || user.Password == "" {
	  c.JSON(400, gin.H{"error": "Username and password cannot be empty"})
	  return
	}
  
	if err := db.Create(&user).Error; err != nil {
	  c.JSON(500, gin.H{"error": "Failed to register user"})
	  return
	}
	c.JSON(200, gin.H{"message": "User registered successfully"})
}

