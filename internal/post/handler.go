package post

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/apperrors"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"net/http"
)

type Handler struct {
	Service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{Service: service}
}

// todo check if this has to be extracted from handlers
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

func buildPostResponse(p *model.Post) *Response {
	return &Response{
		PostID:    p.ID,
		Title:     p.Title,
		Content:   p.Content,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
		Author: UserSummaryResponse{
			UserID:   p.User.ID,
			Username: p.User.Username,
			Email:    p.User.Email,
		},
	}
}

func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("", h.getPosts)
	r.GET("/:id", h.getPost)
	r.POST("", h.createPost)
	r.PATCH("/:id", h.updatePost)
	r.DELETE("/:id", h.deletePost)
}

func (h *Handler) getPosts(c *gin.Context) {
	posts, err := h.Service.GetPosts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := make([]*Response, len(posts))
	for i, p := range posts {
		resp[i] = buildPostResponse(p)
	}
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) getPost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(c, apperrors.NewInvalidInputError("ID must be a uuid"))
		return
	}
	p, err := h.Service.GetPost(id)
	if err != nil {
		handleError(c, err)
	}
	resp := buildPostResponse(p)
	c.JSON(http.StatusOK, resp)

}

func (h *Handler) createPost(c *gin.Context) {
	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewInvalidInputError("Invalid json body"))
		return
	}
	post, err := h.Service.CreatePost(&req)
	if err != nil {
		handleError(c, err)
		return
	}
	resp := buildPostResponse(post)
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) updatePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(c, apperrors.NewInvalidInputError("ID must be a uuid"))
		return
	}
	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleError(c, apperrors.NewInvalidInputError("invalid json body"))
		return
	}
	post, err := h.Service.UpdatePost(id, &req)
	if err != nil {
		handleError(c, err)
		return
	}
	resp := buildPostResponse(post)
	c.JSON(http.StatusOK, resp)
}

func (h *Handler) deletePost(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		handleError(c, apperrors.NewInvalidInputError("ID must be a uuid"))
		return
	}
	err = h.Service.DeletePost(id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
