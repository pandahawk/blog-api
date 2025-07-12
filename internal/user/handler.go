package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", getAllUsers)
	r.GET("/:id", getUserById)
	r.POST("", createUser)
	r.PATCH("/:id", updateUserById)
	r.DELETE("/:id", deleteUserById)
}

func getAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "get all users"})
}

func getUserById(c *gin.Context) {
	id := c.Param("id")
	msg := fmt.Sprintf("get user %s", id)
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func updateUserById(c *gin.Context) {
	id := c.Param("id")
	msg := fmt.Sprintf("update user %s", id)
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func createUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "create new user"})
}

func deleteUserById(c *gin.Context) {
	id := c.Param("id")
	msg := fmt.Sprintf("delete user %s", id)
	c.JSON(http.StatusOK, gin.H{"message": msg})
}
