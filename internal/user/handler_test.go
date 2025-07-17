package user

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/apperrors"
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
	userGroup := router.Group("/users")
	handler.RegisterRoutes(userGroup)
	return router, mockService
}

func TestHandler_GetAllUsers(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		router, mockService := setupTestRouterWithMockService(t)
		wantUser := []User{
			{ID: uuid.New(), Username: "testuser1", Email: "testuser1@example.com"},
			{ID: uuid.New(), Username: "testuser2", Email: "testuser1@example.com"},
		}
		mockService.EXPECT().
			GetAllUsers().
			Return(wantUser, nil)

		req, _ := http.NewRequest(http.MethodGet, "/users", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var gotUser []User
		err := json.NewDecoder(w.Body).Decode(&gotUser)
		assert.NoError(t, err, "Error unmarshaling response body to []User")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.ElementsMatch(t, wantUser, gotUser)
	})
}

func TestHandler_GetUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		router, mockService := setupTestRouterWithMockService(t)
		wantUser := User{ID: uuid.New(), Username: "testuser1"}
		mockService.EXPECT().
			GetUser(gomock.Any()).
			Return(wantUser, nil)

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%v",
			wantUser.ID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var gotUser User
		err := json.NewDecoder(w.Body).Decode(&gotUser)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, wantUser, gotUser)
	})

	t.Run("not found", func(t *testing.T) {
		router, mockService := setupTestRouterWithMockService(t)
		id := uuid.New()
		mockService.EXPECT().
			GetUser(gomock.Any()).
			Return(User{}, apperrors.NewNotFoundError("user", id))

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%v",
			id), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		router, _ := setupTestRouterWithMockService(t)
		id := "abc"

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%v", id), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		router, mockService := setupTestRouterWithMockService(t)
		wantUser := User{
			ID:       uuid.New(),
			Username: "testuser1",
			Email:    "testuser1@example.com",
		}
		mockService.EXPECT().CreateUser(gomock.Any()).Return(wantUser, nil)
		rawJSON := `{"username": "testuser1","email": "testuser1@example.com"}`
		body := strings.NewReader(rawJSON)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", body)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		var gotUser User
		err := json.NewDecoder(w.Body).Decode(&gotUser)
		assert.NoError(t, err)
		assert.Equal(t, wantUser, gotUser)
	})

	t.Run("invalid json", func(t *testing.T) {
		router, _ := setupTestRouterWithMockService(t)
		invalidJSON := `{
	ID:       "1",
	Username: "testuser1",
	Email:    "testuser1@example.com"
	},`
		body := strings.NewReader(invalidJSON)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", body)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("missing required field", func(t *testing.T) {
		router, _ := setupTestRouterWithMockService(t)
		rawJSON := `{
		"id":       1,
		"username": "testuser1",
		"email":    ""
	}`
		body := strings.NewReader(rawJSON)
		w := httptest.NewRecorder()

		req, _ := http.NewRequest(http.MethodPost, "/users", body)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("username taken", func(t *testing.T) {
		router, mockService := setupTestRouterWithMockService(t)
		rawJSON := `{
		"username": "testuser",
		"email":    "testuser@example.com"
	}`
		body := strings.NewReader(rawJSON)
		mockService.EXPECT().CreateUser(gomock.Any()).Return(User{},
			apperrors.NewValidationError("username already exists"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/users", body)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})
}

func TestHandler_UpdateUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		router, mockService := setupTestRouterWithMockService(t)
		id := uuid.New()
		rawJSON := `{
					"email": "testuserupdate@example.com"
				}`
		body := strings.NewReader(rawJSON)
		user := User{
			ID:       id,
			Username: "testuser",
			Email:    "testuserupdate@example.com",
		}
		mockService.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(user, nil)

		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%v",
			id.String()), body)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var gotUser User
		err := json.NewDecoder(w.Body).Decode(&gotUser)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, user, gotUser)
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
		id := 1
		rawJSON := `{
					"username": "updatedtestuser, 
					"email": "testuser@example.com"
				}`
		body := strings.NewReader(rawJSON)

		req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%d", id),
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
		//mockService.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(User{},
		//	apperrors.NewValidationError(""))

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
		mockService.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Return(User{}, apperrors.NewNotFoundError("user", id))

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
		mockService.EXPECT().DeleteUser(gomock.Any()).Return(fmt.Errorf(
			"user %v not found", id))

		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%v",
			id.String()), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), fmt.Sprintf(
			"user %v not found", id.String()))
	})

}
