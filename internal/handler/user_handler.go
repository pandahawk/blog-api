package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/apperrors"
	"github.com/pandahawk/blog-api/internal/dto"
	"github.com/pandahawk/blog-api/internal/mapper"
	"github.com/pandahawk/blog-api/internal/user"
	"log"
	"net/http"
)

type UserHandler struct {
	Service user.Service
}

func NewHandler(service user.Service) *UserHandler {
	return &UserHandler{Service: service}
}

func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{"error": message})
}

func (uh *UserHandler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", uh.getAllUsers)
	r.GET("/:id", uh.getUser)
	r.POST("", uh.createUser)
	r.PATCH("/:id", uh.updateUser)
	r.DELETE("/:id", uh.deleteUser)
}

// @Summary Get all users
// @Description Get all users in the system
// @Tags users
// @Produce json
// @Success 200 {array} user.User
// @Router /users [get]
func (uh *UserHandler) getAllUsers(c *gin.Context) {
	users, _ := uh.Service.GetAllUsers()

	resp := make([]dto.UserResponse, len(users))
	for i, u := range users {
		resp[i] = mapper.FromUser(u)
	}
	c.JSON(http.StatusOK, resp)
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
func (uh *UserHandler) getUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "ID must be a uuid")
		return
	}
	u, err := uh.Service.GetUser(id)
	var ne *apperrors.NotFoundError
	if errors.As(err, &ne) {
		respondWithError(c, http.StatusNotFound, ne.Error())
		return
	}

	resp := mapper.FromUser(u)
	c.JSON(http.StatusOK, resp)
}

// @Summary Create a new user
// @Description Creates a new user and returns the created resource
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} apperrors.ValidationError
// @Failure 409 {object} apperrors.DuplicateError
// @Router /users [post]
func (uh *UserHandler) createUser(c *gin.Context) {
	var req dto.CreateUserRequest
	c.Header("Content-Type", "application/json")
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err.Error())
		respondWithError(c, http.StatusBadRequest, "invalid request body")
		return
	}
	u, err := uh.Service.CreateUser(req)
	var de *apperrors.DuplicateError
	if errors.As(err, &de) {
		respondWithError(c, http.StatusConflict, err.Error())
		return
	}
	resp := mapper.FromUser(u)
	c.JSON(http.StatusCreated, resp)
}

// @Summary Update user by ID
// @Description Updates an existing user and returns the updated resource
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID" format(uuid)
// @Param user body dto.UpdateUserRequest true "User update data"
// @Success 201 {object} user.User
// @Failure 400 {object} apperrors.InvalidInputError
// @Failure 404 {object} apperrors.NotFoundError
// @Router /users/{id} [patch]
func (uh *UserHandler) updateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "ID must be a uuid")
		return
	}
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, err.Error())
		return
	}
	u, err := uh.Service.UpdateUser(id, req)

	var de *apperrors.DuplicateError
	if errors.As(err, &de) {
		respondWithError(c, http.StatusConflict, de.Error())
		return
	}

	var ne *apperrors.NotFoundError
	if errors.As(err, &ne) {
		respondWithError(c, http.StatusNotFound, ne.Error())
		return
	}

	var ie *apperrors.InvalidInputError
	if errors.As(err, &ie) {
		respondWithError(c, http.StatusBadRequest, ie.Error())
	}

	resp := mapper.FromUser(u)
	c.JSON(http.StatusOK, resp)
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
func (uh *UserHandler) deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "ID must be a uuid")
		return
	}

	if err = uh.Service.DeleteUser(id); err != nil {
		respondWithError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
