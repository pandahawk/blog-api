package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/apperrors"
	"github.com/pandahawk/blog-api/internal/dto"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"net/http"
	"strings"
)

type Handler struct {
	Service Service
}

func buildUserResponse(u *model.User) *dto.UserResponse {
	posts := make([]dto.PostSummaryResponse, len(u.Posts))
	for i, p := range u.Posts {
		posts[i] = dto.PostSummaryResponse{
			PostID: p.ID,
			Title:  p.Title,
		}
	}
	return &dto.UserResponse{
		UserID:   u.ID,
		Username: u.Username,
		Email:    u.Email,
		Posts:    posts,
		JoinedAt: u.CreatedAt,
	}
}

func NewHandler(service Service) *Handler {
	return &Handler{Service: service}
}

func handleError(c *gin.Context, err error) {
	var de *apperrors.DuplicateError
	if errors.As(err, &de) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	var ne *apperrors.NotFoundError
	if errors.As(err, &ne) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var ie *apperrors.InvalidInputError
	if errors.As(err, &ie) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.getUsers)
	r.GET("/:id", h.getUser)
	r.POST("", h.createUser)
	r.PATCH("/:id", h.updateUser)
	r.DELETE("/:id", h.deleteUser)
}

// @Summary Get all users
// @Description Get all users in the system
// @Tags users
// @Produce json
// @Success 200 {array} user.User
// @Router /users [get]
func (h *Handler) getUsers(c *gin.Context) {
	users, _ := h.Service.GetUsers()

	resp := make([]*dto.UserResponse, len(users))
	for i, u := range users {
		resp[i] = buildUserResponse(u)
	}
	c.JSON(http.StatusOK, resp)
}

// List all users
// @Summary Get user by ID
// @Description Get the user with the specified ID
// @Tags users
// @Produce json
// @Param id path string true "User ID" format:"uuid"
// @Success 200 {object} user.User
// @Failure 404 {object} apperrors.NotFoundError
// @Failure 400 {object} apperrors.InvalidInputError
// @Router /users/{id} [get]
func (h *Handler) getUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(c, apperrors.NewInvalidInputError("ID must be a uuid"))
		return
	}
	u, err := h.Service.GetUser(id)
	if err != nil {
		handleError(c, err)
		return
	}
	resp := buildUserResponse(u)
	c.JSON(http.StatusOK, resp)
}

// @Summary Create a new user
// @Description Creates a new user and returns the created resource
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User data"
// @Success 201 {object} dto.UserResponse
// @Failure 400 {object} apperrors.InvalidInputError
// @Failure 409 {object} apperrors.DuplicateError
// @Router /users [post]
func (h *Handler) createUser(c *gin.Context) {
	var req dto.CreateUserRequest
	c.Header("Content-Type", "application/json")
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewInvalidInputError("invalid request body"))
		return
	}
	u, err := h.Service.CreateUser(&req)
	if err != nil {
		handleError(c, err)
		return
	}

	resp := buildUserResponse(u)
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
// @Failure 400 {object} apperrors.DuplicateError
// @Router /users/{id} [patch]
func (h *Handler) updateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(c, apperrors.NewInvalidInputError("ID must be a uuid"))
		return
	}
	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if strings.Contains(err.Error(), "Email") {
			handleError(c, apperrors.NewInvalidInputError("invalid email"))
			return
		}
		handleError(c, apperrors.NewInvalidInputError("invalid request body"))
		return
	}
	u, err := h.Service.UpdateUser(id, &req)

	if err != nil {
		handleError(c, err)
		return
	}

	resp := buildUserResponse(u)
	c.JSON(http.StatusOK, resp)
}

// @Summary Delete user by ID
// @Description Deletes an existing user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID" Format(uuid)
// @Success 204
// @Failure 404 {object} apperrors.NotFoundError
// @Failure 400 {object} apperrors.InvalidInputError
// @Router /users/{id} [delete]
func (h *Handler) deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(c, apperrors.NewInvalidInputError("ID must be a uuid"))
		return
	}

	if err = h.Service.DeleteUser(id); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
