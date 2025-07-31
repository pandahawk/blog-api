package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/apperrors"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
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
	userGroup := router.Group("/users")
	handler.RegisterRoutes(userGroup)
	return router, mockService
}

func TestHandler_GetUser(t *testing.T) {
	id := uuid.New()
	tests := []struct {
		name          string
		want          *model.User
		method        string
		path          string
		mockBehaviour func(service *MockService, user *model.User)
		wantStatus    int
	}{
		{
			name: "success",
			want: &model.User{
				ID:       id,
				Username: "testuser",
				Email:    "testuser@example.com",
			},
			method: http.MethodGet,
			path:   fmt.Sprintf("/users/%s", id.String()),
			mockBehaviour: func(service *MockService, user *model.User) {
				service.EXPECT().
					GetUser(gomock.Any()).
					Return(user, nil)
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "not found",
			want:   nil,
			method: http.MethodGet,
			path:   fmt.Sprintf("/users/%s", id.String()),
			mockBehaviour: func(service *MockService, user *model.User) {
				service.EXPECT().
					GetUser(gomock.Any()).
					Return(nil, apperrors.NewNotFoundError("user", id))
			},
			wantStatus: http.StatusNotFound,
		},
		{
			name:          "invalid id",
			want:          nil,
			method:        http.MethodGet,
			path:          "/users/abc",
			mockBehaviour: nil,
			wantStatus:    http.StatusBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, mockService := setupTestRouterWithMockService(t)
			if test.mockBehaviour != nil {
				test.mockBehaviour(mockService, test.want)
			}

			w := httptest.NewRecorder()
			req, _ := http.NewRequest(test.method, test.path, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, test.wantStatus, w.Code)
		})
	}
}

func TestHandler_GetUsers(t *testing.T) {
	tests := []struct {
		name          string
		mockBehaviour func(service *MockService, users []*model.User)
		wantUsers     []*model.User
		want          []*Response
		wantStatus    int
		wantErr       string
	}{
		{
			name: "success",
			mockBehaviour: func(service *MockService, users []*model.User) {
				service.EXPECT().GetUsers().Return(users, nil)
			},
			wantUsers: []*model.User{
				{ID: uuid.New(), Username: "testuser1", Email: "testuser1@mail.com"},
				{ID: uuid.New(), Username: "testuser2", Email: "testuser2@mail.com"},
			},
			want: []*Response{
				{Username: "testuser1", Email: "testuser1@mail.com"},
				{Username: "testuser2", Email: "testuser2@mail.com"},
			},
			wantStatus: 200,
			wantErr:    "",
		},
		{
			name:      "db error",
			wantUsers: nil,
			want:      nil,
			mockBehaviour: func(service *MockService, users []*model.User) {
				service.EXPECT().GetUsers().
					Return(nil, errors.New("db error"))
			},
			wantStatus: 500,
			wantErr:    "db error",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, mockService := setupTestRouterWithMockService(t)
			if test.mockBehaviour != nil {
				test.mockBehaviour(mockService, test.wantUsers)
			}
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/users", nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, test.wantStatus, w.Code)
			if test.wantErr == "" {
				var r []*Response
				if err := json.NewDecoder(w.Body).Decode(&r); err != nil {
					log.Fatal(err)
				}
				for i, expected := range test.want {
					assert.Equal(t, expected.Username, r[i].Username)
					assert.Equal(t, expected.Email, r[i].Email)
				}
			}
		})
	}
}

func TestHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		rawBody       string
		mockBehaviour func(service *MockService)
		wantStatus    int
		want          *Response
		wantErr       string
	}{
		{
			name:    "success",
			rawBody: `{"username": "testuser", "email":"testuser@example.com"}`,
			mockBehaviour: func(service *MockService) {
				service.EXPECT().CreateUser(gomock.Any()).
					Return(&model.User{
						ID:       uuid.Nil,
						Username: "testuser",
						Email:    "testuser@mail.com",
					}, nil)
			},
			want: &Response{
				UserID:   uuid.Nil,
				Username: "testuser",
				Email:    "testuser@mail.com",
			},
			wantStatus: http.StatusCreated,
			wantErr:    "",
		},
		{
			name:          "invalid json",
			rawBody:       `{"username": "testuser", "email":"testuser@example.com"`,
			mockBehaviour: nil,
			want:          nil,
			wantStatus:    http.StatusBadRequest,
			wantErr:       "invalid request body",
		},
		{
			name:    "invalid username",
			rawBody: `{"username":"123","email":"testuser@mail.com"}`,
			mockBehaviour: func(service *MockService) {
				service.EXPECT().CreateUser(gomock.Any()).
					Return(nil, apperrors.NewInvalidInputError(
						"invalid username: must not be a number"))
			},
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    "invalid username: must not be a number",
		},
		{
			name:    "duplicate username",
			rawBody: `{"username":"testuser","email":"testuser@mail.com"}`,
			mockBehaviour: func(service *MockService) {
				service.EXPECT().CreateUser(gomock.Any()).
					Return(nil, apperrors.NewDuplicateError("username"))
			},
			want:       nil,
			wantStatus: http.StatusConflict,
			wantErr:    "username already exists",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, mockService := setupTestRouterWithMockService(t)
			if test.mockBehaviour != nil {
				test.mockBehaviour(mockService)
			}
			body := strings.NewReader(test.rawBody)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/users", body)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, test.wantStatus, w.Code)
			if test.wantStatus == http.StatusCreated {
				var r *Response
				err := json.NewDecoder(w.Body).Decode(&r)
				if err != nil {
					log.Fatal(err)
				}
				assert.NoError(t, err)
				assert.Equal(t, test.want.Username, r.Username)
				assert.Equal(t, test.want.Email, r.Email)
			} else {
				bodyBytes, _ := io.ReadAll(w.Body)
				bodyStr := string(bodyBytes)
				assert.Contains(t, bodyStr, test.wantErr)
			}
		})
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		rawBody       string
		wantUser      *model.User
		mockBehaviour func(service *MockService, id string, rawBody string, user *model.User)
		wantStatus    int
		wantResponse  *Response
		wantErr       string
	}{
		{
			name:    "success",
			id:      uuid.Nil.String(),
			rawBody: `{"username":"updated","email":"updated@mail.com"}`,
			wantUser: &model.User{
				ID:       uuid.Nil,
				Username: "updated",
				Email:    "updated@mail.com",
			},
			mockBehaviour: func(service *MockService, id string,
				rawBody string, user *model.User) {
				service.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).
					Return(user, nil)
			},
			wantStatus: http.StatusOK,
			wantResponse: &Response{
				UserID:   uuid.Nil,
				Username: "updated",
				Email:    "updated@mail.com"},
			wantErr: "",
		},
		{
			name:          "not an uuid",
			id:            "abc",
			rawBody:       `{"username": "updated", "email":"updated@mail.com"}`,
			wantUser:      nil,
			mockBehaviour: nil,
			wantStatus:    http.StatusBadRequest,
			wantResponse:  nil,
			wantErr:       "ID must be a uuid",
		},
		{
			name:          "invalid json",
			id:            uuid.Nil.String(),
			rawBody:       `{"username": "updated", "email":"updated@mail.com"`,
			wantUser:      nil,
			mockBehaviour: nil,
			wantStatus:    http.StatusBadRequest,
			wantResponse:  nil,
			wantErr:       "invalid request body",
		},
		{
			name:          "invalid email",
			id:            uuid.Nil.String(),
			rawBody:       `{"username": "updated", "email":"updatedmail.com"}`,
			wantUser:      nil,
			mockBehaviour: nil,
			wantStatus:    http.StatusBadRequest,
			wantResponse:  nil,
			wantErr:       "invalid email",
		},
		{
			name:     "invalid username",
			id:       uuid.Nil.String(),
			rawBody:  `{"username": "123", "email":"updated@mail.com"}`,
			wantUser: nil,
			mockBehaviour: func(service *MockService, id string,
				rawBody string, user *model.User) {
				service.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).
					Return(nil, apperrors.NewInvalidInputError(
						"invalid username: must not be a number"))
			},
			wantStatus:   http.StatusBadRequest,
			wantResponse: nil,
			wantErr:      "invalid username: must not be a number",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, mockService := setupTestRouterWithMockService(t)
			if test.mockBehaviour != nil {
				test.mockBehaviour(mockService, test.id, test.rawBody, test.wantUser)
			}
			body := strings.NewReader(test.rawBody)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPatch,
				"/users/"+test.id, body)
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			assert.Equal(t, test.wantStatus, w.Code)
			if test.wantErr == "" {
				var r *Response
				if err := json.NewDecoder(w.Body).Decode(&r); err != nil {
					log.Fatal(err)
				}
				assert.Equal(t, test.wantUser.Username, r.Username)
				assert.Equal(t, test.wantUser.Email, r.Email)
				assert.Equal(t, test.wantUser.ID, r.UserID)
			} else {
				bodyBytes, _ := io.ReadAll(w.Body)
				bodyStr := string(bodyBytes)
				assert.Contains(t, bodyStr, test.wantErr)
			}
		})
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		wantStatus    int
		wantErr       string
		mockBehaviour func(service *MockService, id string)
	}{
		{
			name:       "success",
			id:         uuid.Nil.String(),
			wantStatus: 204,
			wantErr:    "",
			mockBehaviour: func(service *MockService, id string) {
				service.EXPECT().DeleteUser(gomock.Any()).
					Return(nil)
			},
		},
		{
			name:          "not an uuid",
			id:            "abc",
			wantStatus:    400,
			wantErr:       "ID must be a uuid",
			mockBehaviour: nil,
		},
		{
			name:       "failed",
			id:         uuid.Nil.String(),
			wantStatus: 500,
			wantErr:    "deletion failed",
			mockBehaviour: func(service *MockService, id string) {
				service.EXPECT().DeleteUser(gomock.Any()).
					Return(errors.New("deletion failed"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, mockService := setupTestRouterWithMockService(t)
			if test.mockBehaviour != nil {
				test.mockBehaviour(mockService, test.id)
			}
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodDelete,
				"/users/"+test.id, nil)
			router.ServeHTTP(w, req)

			assert.Equal(t, test.wantStatus, w.Code)
			if test.wantErr != "" {
				assert.Contains(t, w.Body.String(), test.wantErr)
			}
		})
	}
}

func TestHandler_buildUserResponse(t *testing.T) {
	id := uuid.New()
	posts := []*model.Post{
		model.NewPost("title1", "content1", id),
		model.NewPost("title2", "content2", id),
	}
	user := &model.User{
		ID:       id,
		Username: "testuser",
		Email:    "testuser@mail.com",
		Posts:    posts,
	}
	want := &Response{
		UserID:   id,
		Username: "testuser",
		Email:    "testuser@mail.com",
		Posts: []*PostSummaryResponse{
			{
				PostID: posts[0].ID,
				Title:  "title1",
			},
			{
				PostID: posts[1].ID,
				Title:  "title2",
			}},
	}

	got := buildUserResponse(user)
	assert.Equal(t, want, got)
}
