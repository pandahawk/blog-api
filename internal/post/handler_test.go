package post

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/apperrors"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
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
		wantPosts     []*model.Post
		want          []*Response
		mockbehaviour func(service *MockService, posts []*model.Post)
		wantStatus    int
		wantErr       string
	}{
		{
			name: "success",
			wantPosts: []*model.Post{
				{
					ID:      uuid.New(),
					Title:   "title1",
					Content: "content1",
					User: &model.User{
						ID:       uuid.New(),
						Username: "user1",
						Email:    "user1@mail.com"},
				},
				{
					ID:      uuid.New(),
					Title:   "title2",
					Content: "content2",
					User: &model.User{
						ID:       uuid.New(),
						Username: "user2",
						Email:    "user2@mail.com"}},
			},
			want: []*Response{
				{
					Title:   "title1",
					Content: "content1",
					Author: UserSummaryResponse{
						UserID:   uuid.Nil,
						Username: "user1",
					},
				},
				{
					Title:   "title2",
					Content: "content2",
					Author: UserSummaryResponse{
						UserID:   uuid.Nil,
						Username: "user2",
					},
				},
			},
			mockbehaviour: func(service *MockService, posts []*model.Post) {
				service.EXPECT().GetPosts().Return(posts, nil)
			},
			wantStatus: 200,
			wantErr:    "",
		},
		{
			name:      "failed",
			wantPosts: nil,
			want:      nil,
			mockbehaviour: func(service *MockService, posts []*model.Post) {
				service.EXPECT().GetPosts().Return(nil, errors.New("failed to get posts"))
			},
			wantStatus: 500,
			wantErr:    "failed to get posts",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, mockService := setupTestRouterWithMockService(t)
			if test.mockbehaviour != nil {
				test.mockbehaviour(mockService, test.wantPosts)

				w := httptest.NewRecorder()
				req, _ := http.NewRequest(http.MethodGet, "/posts", nil)
				router.ServeHTTP(w, req)

				assert.Equal(t, w.Code, test.wantStatus)
				if test.wantErr == "" {
					var r []*Response
					if err := json.NewDecoder(w.Body).Decode(&r); err != nil {
						t.Fatal(err)
					}
					for i, response := range r {
						assert.Equal(t, response, r[i])
						assert.Equal(t, test.wantPosts[i].Title, response.Title)
						assert.Equal(t, test.wantPosts[i].User.Username, response.Author.Username)
						assert.Equal(t, test.wantPosts[i].Content, response.Content)
					}
				} else {
					assert.Contains(t, w.Body.String(), test.wantErr)
				}
			}
		})
	}
}

func TestHandler_GetPost(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		wantPost      *model.Post
		mockbehaviour func(service *MockService, id string, post *model.Post)
		wantStatus    int
		wantErr       string
	}{
		{
			name: "success",
			id:   uuid.Nil.String(),
			wantPost: &model.Post{
				ID:      uuid.Nil,
				Title:   "title",
				Content: "content",
				UserID:  uuid.New(),
				User: &model.User{
					ID:       uuid.New(),
					Username: "user1",
					Email:    "user1@mail.com",
				},
			},
			mockbehaviour: func(service *MockService, id string, post *model.Post) {
				service.EXPECT().GetPost(gomock.Any()).Return(post, nil)
			},
			wantStatus: 200,
			wantErr:    "",
		},
		{
			name:     "not found",
			id:       uuid.Nil.String(),
			wantPost: nil,
			mockbehaviour: func(service *MockService, id string, post *model.Post) {
				service.EXPECT().GetPost(gomock.Any()).
					Return(nil, apperrors.NewNotFoundError("post", uuid.Nil))
			},
			wantStatus: 404,
			wantErr:    "not found",
		},
		{
			name: "not an uuid",
			id:   "abc",
			wantPost: &model.Post{
				ID:      uuid.Nil,
				Title:   "title",
				Content: "content",
				UserID:  uuid.Nil,
				User:    nil,
			},
			mockbehaviour: nil,
			wantStatus:    400,
			wantErr:       "ID must be a uuid",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, mockService := setupTestRouterWithMockService(t)
			if test.mockbehaviour != nil {
				test.mockbehaviour(mockService, test.id, test.wantPost)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/posts/"+test.id, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.wantStatus)
			if test.wantErr == "" {

			} else {
				assert.Contains(t, w.Body.String(), test.wantErr)
			}

		})
	}
}

func TestHandler_CreatePost(t *testing.T) {
	tests := []struct {
		name          string
		rawBody       string
		wantPost      *model.Post
		mockBehaviour func(service *MockService, post *model.Post)
		wantStatus    int
		wantErr       string
	}{
		{
			name: "success",
			rawBody: `
				{"title":"title",
				"content":"content",
				"author_id":"3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a"}
				`,
			wantPost: &model.Post{
				ID:      uuid.MustParse("3c9a5f8d-91c6-4e3e-9f76-046b7e9b6c1a"),
				Title:   "title",
				Content: "content",
				UserID:  uuid.MustParse("3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a"),
				User: &model.User{ID: uuid.MustParse(
					"3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a")},
			},
			mockBehaviour: func(service *MockService, post *model.Post) {
				service.EXPECT().CreatePost(gomock.Any()).Return(post, nil)
			},
			wantStatus: 201,
			wantErr:    "",
		},
		{
			name: "invalid json body",
			rawBody: `
				{"title":"title",
				"content":"content",
				"author_id":"3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a"
				`,
			wantPost:      nil,
			mockBehaviour: nil,
			wantStatus:    400,
			wantErr:       "invalid json",
		},
		{
			name: "invalid title",
			rawBody: `
				{"title":"123",
				"content":"content",
				"author_id":"3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a"}
				`,
			wantPost: nil,
			mockBehaviour: func(service *MockService, post *model.Post) {
				service.EXPECT().CreatePost(gomock.Any()).
					Return(nil, apperrors.NewInvalidInputError(
						"invalid title"))
			},
			wantStatus: 400,
			wantErr:    "invalid",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, mockService := setupTestRouterWithMockService(t)
			if test.mockBehaviour != nil {
				test.mockBehaviour(mockService, test.wantPost)
			}

			w := httptest.NewRecorder()
			body := strings.NewReader(test.rawBody)
			req, _ := http.NewRequest(http.MethodPost, "/posts", body)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, test.wantStatus, w.Code)
			if test.wantErr == "" {
				var r *Response
				if err := json.NewDecoder(w.Body).Decode(&r); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, test.wantPost.Title, r.Title)
				assert.Equal(t, test.wantPost.Content, r.Content)
			}
		})
	}
}

func TestHandler_UpdatePost(t *testing.T) {
	tests := []struct {
		name          string
		rawBody       string
		id            string
		wantPost      *model.Post
		mockBehaviour func(service *MockService, id uuid.UUID, post *model.Post)
		wantStatus    int
		wantErr       string
	}{
		{
			name: "success",
			rawBody: `
				{"title":"updated",
				"content":"content updated",
				"author_id":"3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a"}
				`,
			id: "3c9a5f8d-91c6-4e3e-9f76-046b7e9b6c1a",
			wantPost: &model.Post{
				ID:      uuid.MustParse("3c9a5f8d-91c6-4e3e-9f76-046b7e9b6c1a"),
				Title:   "updated",
				Content: "content updated",
				UserID:  uuid.MustParse("3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a"),
				User: &model.User{ID: uuid.MustParse(
					"3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a")},
			},
			mockBehaviour: func(service *MockService, id uuid.UUID, post *model.Post) {
				service.EXPECT().UpdatePost(gomock.Any(), gomock.Any()).Return(post, nil)
			},
			wantStatus: 200,
			wantErr:    "",
		},
		{
			name: "not an uuid",
			rawBody: `
				{"title":"updated",
				"content":"content updated",
				"author_id":"3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a"}
				`,
			id:            "abc",
			wantPost:      nil,
			mockBehaviour: nil,
			wantStatus:    400,
			wantErr:       "must be a uuid",
		},
		{
			name: "invalid json body",
			rawBody: `
				{"title":"updated",
				"content":"content updated",
				"author_id":"3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a"
				`,
			id:            "3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a",
			wantPost:      nil,
			mockBehaviour: nil,
			wantStatus:    400,
			wantErr:       "invalid json",
		},
		{
			name: "blank content",
			rawBody: `
					{"title":"updated",
					"content":" ",
					"author_id":"3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a"}
				`,
			id: "3c9a5f8d-91c6-4e3e-9f76-046b7e9b6c1a",
			wantPost: &model.Post{
				ID:      uuid.MustParse("3c9a5f8d-91c6-4e3e-9f76-046b7e9b6c1a"),
				Title:   "updated",
				Content: "content updated",
				UserID:  uuid.MustParse("3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a"),
				User: &model.User{ID: uuid.MustParse(
					"3c9d5f8d-91c6-4e3e-9f76-046b7e9b6c1a")},
			},
			mockBehaviour: func(service *MockService, id uuid.UUID, post *model.Post) {
				service.EXPECT().UpdatePost(gomock.Any(), gomock.Any()).
					Return(nil, apperrors.NewInvalidInputError("content must not be blank"))
			},
			wantStatus: 400,
			wantErr:    "not be blank",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, mockService := setupTestRouterWithMockService(t)
			if test.mockBehaviour != nil {
				test.mockBehaviour(mockService, uuid.MustParse(test.id), test.wantPost)
			}

			w := httptest.NewRecorder()
			body := strings.NewReader(test.rawBody)
			req, _ := http.NewRequest(http.MethodPatch, "/posts/"+test.id, body)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, test.wantStatus, w.Code)
			if test.wantErr == "" {
				var r Response
				if err := json.NewDecoder(w.Body).Decode(&r); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, test.wantPost.Title, r.Title)
				assert.Equal(t, test.wantPost.Content, r.Content)
			} else {
				assert.Contains(t, w.Body.String(), test.wantErr)
			}
		})
	}
}

func TestHandler_DeletePost(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		mockBehaviour func(service *MockService, id uuid.UUID)
		wantStatus    int
		wantErr       string
	}{
		{
			name: "success",
			id:   uuid.Nil.String(),
			mockBehaviour: func(service *MockService, id uuid.UUID) {
				service.EXPECT().DeletePost(id).Return(nil)
			},
			wantStatus: 204,
			wantErr:    "",
		},
		{
			name:          "not a uuid",
			id:            "abc",
			mockBehaviour: nil,
			wantStatus:    400,
			wantErr:       "must be a uuid",
		},
		{
			name: "post not found",
			id:   uuid.Nil.String(),
			mockBehaviour: func(service *MockService, id uuid.UUID) {
				service.EXPECT().DeletePost(id).
					Return(apperrors.NewNotFoundError("post", uuid.Nil))
			},
			wantStatus: 404,
			wantErr:    "not found",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, mockService := setupTestRouterWithMockService(t)
			if test.mockBehaviour != nil {
				test.mockBehaviour(mockService, uuid.MustParse(test.id))
			}
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete, "/posts/"+test.id, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, test.wantStatus, w.Code)
			if test.wantErr == "" {
				assert.Empty(t, w.Body.String())
			} else {
				assert.Contains(t, w.Body.String(), test.wantErr)
			}

		})
	}
}
