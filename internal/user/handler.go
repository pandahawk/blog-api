package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Handler struct {
	Service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{Service: service}
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.getAllUsers)
	r.GET("/:id", h.getUser)
	r.POST("", h.createUser)
	r.PATCH("/:id", h.updateUser)
	r.DELETE("/:id", h.deleteUser)
}

func (h *Handler) getAllUsers(c *gin.Context) {
	users := h.Service.GetAllUsers()
	c.JSON(http.StatusOK, users)
}

func (h *Handler) getUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid ID")
		return
	}
	user, err := h.Service.GetUser(id)
	if err != nil {
		respondWithError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) createUser(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	savedUser, err := h.Service.CreateUser(newUser)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, savedUser)
}

func (h *Handler) updateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid ID")
		return
	}
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	updatedUser, err := h.Service.UpdateUser(id, newUser)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func (h *Handler) deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "invalid ID")
		return
	}
	err = h.Service.DeleteUser(id)
	if err != nil {
		respondWithError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
