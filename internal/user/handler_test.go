package user

import (
	"bytes"
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
	"time"
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

func TestHandler_GetUsersOld(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		router, mockService := setupTestRouterWithMockService(t)
		wantUser := []*model.User{
			{ID: uuid.New(), Username: "testuser1",
				Email: "testuser1@example.com", CreatedAt: time.Now()},
			{ID: uuid.New(), Username: "testuser2",
				Email: "testuser1@example.com", CreatedAt: time.Now()},
		}
		mockService.EXPECT().
			GetUsers().
			Return(wantUser, nil)

		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var ur []Response
		err := json.NewDecoder(w.Body).Decode(&ur)
		assert.NoError(t, err, "Error unmarshaling response body to []User")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, wantUser[0].ID, ur[0].UserID)
		assert.Equal(t, wantUser[1].ID, ur[1].UserID)
	})
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

// TODO refactor this and other tests
func TestHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name          string
		req           *CreateUserRequest
		rawBody       string
		mockBehaviour func(service *MockService)
		wantStatus    int
		want          *Response
		wantErr       string
	}{
		{
			name:    "success",
			rawBody: "",
			req: &CreateUserRequest{
				Username: "testuser",
				Email:    "testuser@mail.com",
			},
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
			rawBody: "",
			req: &CreateUserRequest{
				Username: "123",
				Email:    "testuser@mail.com",
			},
			mockBehaviour: func(service *MockService) {
				service.EXPECT().CreateUser(gomock.Any()).
					Return(nil, apperrors.NewInvalidInputError(
						"invalid username: must not be a number"))
			},
			want:       nil,
			wantStatus: http.StatusBadRequest,
			wantErr:    "\"invalid username: must not be a number\"",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router, mockService := setupTestRouterWithMockService(t)
			if test.mockBehaviour != nil {
				test.mockBehaviour(mockService)
			}

			var body io.Reader
			if test.rawBody != "" {
				body = strings.NewReader(test.rawBody)
			} else {
				buf := new(bytes.Buffer)
				if err := json.NewEncoder(buf).Encode(test.req); err != nil {
					log.Fatal(err)
				}
				body = buf
			}

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

	t.Run("success", func(t *testing.T) {
		router, mockService := setupTestRouterWithMockService(t)
		id := uuid.New()
		rawJSON := `{
					"email": "testuserupdate@example.com"
				}`
		body := strings.NewReader(rawJSON)
		wantUser := model.User{
			ID:        id,
			Username:  "testuser",
			Email:     "testuserupdate@example.com",
			CreatedAt: time.Now(),
		}
		mockService.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(&wantUser, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%v",
			id.String()), body)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		var ur Response
		err := json.NewDecoder(w.Body).Decode(&ur)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, wantUser.ID, ur.UserID)
		assert.Equal(t, wantUser.Email, ur.Email)
		assert.Equal(t, wantUser.Username, ur.Username)

	})
	t.Run("email is not valid", func(t *testing.T) {
		router, _ := setupTestRouterWithMockService(t)
		id := uuid.New()
		rawJSON := `{
					"email": "testuserupdateexample.com"
				}`
		body := strings.NewReader(rawJSON)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%v",
			id.String()), body)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("invalid id", func(t *testing.T) {
		router, _ := setupTestRouterWithMockService(t)
		id := "abc"
		rawJSON := `{
					"username": "updatedtestuser", 
					"email": "testuser@example.com"
				}`
		body := strings.NewReader(rawJSON)

		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%s", id),
			body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("bad json", func(t *testing.T) {
		router, _ := setupTestRouterWithMockService(t)
		id := uuid.New()
		rawJSON := `{
					"username": "updatedtestuser, 
					"email": "testuser@example.com"
				}`
		body := strings.NewReader(rawJSON)

		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%s", id),
			body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("empty fields", func(t *testing.T) {
		router, _ := setupTestRouterWithMockService(t)
		id := 1
		rawJSON := `{
					"username": "", 
					"email": ""
				}`
		body := strings.NewReader(rawJSON)

		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%d", id), body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
	t.Run("not found", func(t *testing.T) {
		router, mockService := setupTestRouterWithMockService(t)
		id := uuid.New()
		rawJSON := `{
					"email": "testuserupdate@example.com"
				}`
		body := strings.NewReader(rawJSON)
		mockService.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).
			Return(nil, apperrors.NewNotFoundError("user", id))

		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%v",
			id.String()), body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHandler_DeleteUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		router, mockService := setupTestRouterWithMockService(t)
		id := uuid.New()
		mockService.EXPECT().DeleteUser(gomock.Any()).Return(nil)

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s",
			id.String()), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		router, _ := setupTestRouterWithMockService(t)
		id := "not a uuid"

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s",
			id), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		router, mockService := setupTestRouterWithMockService(t)
		id := uuid.New()
		mockService.EXPECT().DeleteUser(gomock.Any()).
			Return(apperrors.NewNotFoundError("user", id))

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%v",
			id.String()), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(),
			apperrors.NewNotFoundError("user", id).Error())
	})

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
