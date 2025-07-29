package post

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"testing"
)

func setupTestRouterWithMockService(t *testing.T) (*gin.Engine, *MockService) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockService := NewMockService(ctrl)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := NewHandler(mockService)
	userGroup := router.Group("/posts")
	handler.RegisterRoutes(userGroup)
	return router, mockService
}

func TestHandler_GetPosts(t *testing.T) {
	tests := []struct {
		name          string
		wantPost      *model.Post
		want          *Response
		mockbehaviour func(service *MockService, post *model.Post)
	}{
		{},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//router, mockService := setupTestRouterWithMockService(t)
			//if test.mockbehaviour != nil {
			//	test.mockbehaviour(mockService, test.wantPost)
			//}
		})
	}
}
