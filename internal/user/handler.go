package user

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	Service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.getAllUsers)
	r.GET("/:id", h.getUser)
	r.POST("", h.createUser)
}

func (h *Handler) getAllUsers(c *gin.Context) {
	users := h.Service.GetAllUsers()
	c.IndentedJSON(http.StatusOK, users)
}

func (h *Handler) getUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	user, err := h.Service.GetUser(id)
	if err != nil {
		log.Printf("failed to get user by ID %v: %v", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusOK, user)
}

func (h *Handler) createUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	savedUser, err := h.Service.CreateUser(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.IndentedJSON(http.StatusCreated, savedUser)
}
