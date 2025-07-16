package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/pandahawk/blog-api/internal/apperrors"
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
	c.JSON(code, gin.H{"errors": message})
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.getAllUsers)
	r.GET("/:id", h.getUser)
	r.POST("", h.createUser)
	r.PATCH("/:id", h.updateUser)
	r.DELETE("/:id", h.deleteUser)
}

func (h *Handler) getAllUsers(c *gin.Context) {
	users, _ := h.Service.GetAllUsers()
	c.JSON(http.StatusOK, users)
}

func (h *Handler) getUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "ID must be an integer")
		return
	}
	user, err := h.Service.GetUser(id)
	var ne *apperrors.NotFoundError
	if errors.As(err, &ne) {
		respondWithError(c, http.StatusNotFound, ne.Error())
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *Handler) createUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "missing username or email")
		return
	}
	savedUser, err := h.Service.CreateUser(req)
	var ve *apperrors.ValidationError
	if errors.As(err, &ve) {
		respondWithError(c, http.StatusConflict, err.Error())
		return
	}
	c.JSON(http.StatusCreated, savedUser)
}

func (h *Handler) updateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "ID must be an integer")
		return
	}
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	updatedUser, err := h.Service.UpdateUser(id, req)

	var ve *apperrors.ValidationError
	if errors.As(err, &ve) {
		respondWithError(c, http.StatusBadRequest, ve.Error())
		return
	}

	var ne *apperrors.NotFoundError
	if errors.As(err, &ne) {
		respondWithError(c, http.StatusNotFound, ne.Error())
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func (h *Handler) deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "ID must be an integer")
		return
	}

	if err = h.Service.DeleteUser(id); err != nil {
		respondWithError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
