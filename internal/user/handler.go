package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	Service UserService
}

func NewHandler(service UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.getAllUsers)
	r.GET("/:id", h.getUserById)
	r.POST("", h.createUser)
	r.PATCH("/:id", h.updateUserById)
	r.DELETE("/:id", h.deleteUserById)
}

func (h *UserHandler) getAllUsers(c *gin.Context) {
	msg := h.Service.GetAllUsers()
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (h *UserHandler) getUserById(c *gin.Context) {
	id := c.Param("id")
	msg := h.Service.GetUser(id)
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (h *UserHandler) updateUserById(c *gin.Context) {
	id := c.Param("id")
	msg := h.Service.UpdateUser(id)
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (h *UserHandler) createUser(c *gin.Context) {
	msg := h.Service.CreateUser()
	c.JSON(http.StatusCreated, gin.H{"message": msg})
}

func (h *UserHandler) deleteUserById(c *gin.Context) {
	id := c.Param("id")
	msg := h.Service.DeleteUser(id)
	c.JSON(http.StatusOK, gin.H{"message": msg})
}
