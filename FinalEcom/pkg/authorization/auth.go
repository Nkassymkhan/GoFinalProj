package authorization

import (
	"github.com/gin-gonic/gin"
)
func Register(c *gin.Context) {
	var user User
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