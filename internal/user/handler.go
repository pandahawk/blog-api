package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/apperrors"
	"github.com/pandahawk/blog-api/internal/dto"
	"log"
	"net/http"
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

// @Summary Get all users
// @Description Get all users in the system
// @Tags users
// @Produce json
// @Success 200 {array} User
// @Router /users [get]
func (h *Handler) getAllUsers(c *gin.Context) {
	users, _ := h.Service.GetAllUsers()

	c.JSON(http.StatusOK, users)
}

// List all users
// @Summary Get user by ID
// @Description Get the user with the specified ID
// @Tags users
// @Produce json
// @Param id path string true "User ID" Format(uuid)
// @Success 200 {object} user.User
// @Failure 404 {object} apperrors.NotFoundError
// @Router /users/{id} [get]
func (h *Handler) getUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
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
	resp := dto.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Posts:    nil,
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary Create a new user
// @Description Creates a new user and returns the created resource
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.CreateUserRequest true "User data"
// @Success 201 {object} user.User
// @Failure 400 {object} apperrors.ValidationError
// @Failure 409 {object} apperrors.ValidationError
// @Router /users [post]
func (h *Handler) createUser(c *gin.Context) {
	var req dto.CreateUserRequest
	c.Header("Content-Type", "application/json")
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error())
		respondWithError(c, http.StatusBadRequest, "uername and email are required")
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

// @Summary Update user by ID
// @Description Updates an existing user and returns the updated resource
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID" Format(uuid)
// @Param user body user.UpdateUserRequest true "User update data"
// @Success 201 {object} user.User
// @Failure 400 {object} apperrors.ValidationError
// @Failure 409 {object} apperrors.ValidationError
// @Failure 404 {object} apperrors.NotFoundError
// @Router /users/{id} [patch]
func (h *Handler) updateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "ID must be an integer")
		return
	}
	var req dto.UpdateUserRequest
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

// @Summary Delete user by ID
// @Description Deletes an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID" Format(uuid)
// @Success 204
// @Failure 400 {object} apperrors.ValidationError
// @Failure 404 {object} apperrors.NotFoundError
// @Router /users/{id} [delete]
func (h *Handler) deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
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
