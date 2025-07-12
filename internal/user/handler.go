package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Handler struct {
	Service Service
}

//	func NewHandler(service Service) *Handler {
//		return &Handler{Service: service}
//	}
func NewHandler(service Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.getAllUsers)
	r.GET("/:id", h.getUser)
	//r.POST("", h.createUser)
	//r.PATCH("/:id", h.updateUserById)
	//r.DELETE("/:id", h.deleteUserById)
}

func (h *Handler) getAllUsers(c *gin.Context) {
	users := h.Service.GetAllUsers()
	c.IndentedJSON(http.StatusOK, users)
}

func (h *Handler) getUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.Service.GetUser(id)
	if err != nil {
		log.Printf("failed to get user by ID %s: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

//
//func (h *Handler) updateUserById(c *gin.Context) {
//	id := c.Param("id")
//	msg := h.Service.UpdateUser(id)
//	c.JSON(http.StatusOK, gin.H{"message": msg})
//}
//
//func (h *Handler) createUser(c *gin.Context) {
//	msg := h.Service.CreateUser()
//	c.JSON(http.StatusCreated, gin.H{"message": msg})
//}
//
//func (h *Handler) deleteUserById(c *gin.Context) {
//	id := c.Param("id")
//	msg := h.Service.DeleteUser(id)
//	c.JSON(http.StatusOK, gin.H{"message": msg})
//}
