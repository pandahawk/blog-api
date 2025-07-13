package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestRouter(service Service) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := NewHandler(service)
	userGroup := router.Group("/users")
	handler.RegisterRoutes(userGroup)
	return router
}

func Test_getAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockService(ctrl)

	want := []User{
		{ID: 1, Username: "testuser1", Email: "testuser1@example.com"},
		{ID: 2, Username: "testuser2", Email: "testuser1@example.com"},
	}
	mockService.EXPECT().
		GetAllUsers().
		Return(want)

	router := setupTestRouter(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	var got []User
	err := json.NewDecoder(w.Body).Decode(&got)

	assert.NoError(t, err, "Error unmarshaling response body to []User")
	assert.Equal(t, http.StatusOK, w.Code)
	assert.ElementsMatch(t, want, got)
}

func Test_getUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockService(ctrl)

	id := 1
	want := User{ID: id, Username: "testuser1"}

	mockService.EXPECT().
		GetUser(gomock.Any()).
		Return(want, nil)
	router := setupTestRouter(mockService)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%v", id), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var got User
	err := json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err, "Error unmarshaling response body to User")

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, want, got)
}

func Test_getUser_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockService(ctrl)

	id := 100

	mockService.EXPECT().
		GetUser(gomock.Any()).
		Return(User{}, fmt.Errorf("user %d not found", id))
	router := setupTestRouter(mockService)

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%d", id), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), fmt.Sprintf("user %d not found", id))
}

func Test_getUser_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockService(ctrl)

	id := "abc"
	router := setupTestRouter(mockService)
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%v", id), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), `"error":"invalid ID"`)
}

func Test_createUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockService(ctrl)
	user := User{
		ID:       1,
		Username: "testuser1",
		Email:    "testuser1@example.com",
	}

	mockService.EXPECT().CreateUser(gomock.Any()).Return(user, nil)
	bodyBytes, _ := json.Marshal(user)
	w := httptest.NewRecorder()

	router := setupTestRouter(mockService)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
	router.ServeHTTP(w, req)

	var got User
	err := json.NewDecoder(w.Body).Decode(&got)
	assert.NoError(t, err)
	assert.Equal(t, user, got)
}

func Test_createUser_invalidJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockService(ctrl)
	invalidJSON := `{
	ID:       "1",
	Username: "testuser1",
	Email:    "testuser1@example.com"
	},`

	w := httptest.NewRecorder()

	router := setupTestRouter(mockService)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Invalid status code: %d but expected 400", w.Code)
	}
}

func Test_createUserWithoutEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockService(ctrl)
	user := User{
		ID:       1,
		Username: "testuser1",
		Email:    "",
	}

	mockService.EXPECT().CreateUser(gomock.Any()).Return(User{},
		errors.New("missing required fields"))
	bodyBytes, _ := json.Marshal(user)
	w := httptest.NewRecorder()

	router := setupTestRouter(mockService)
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewReader(bodyBytes))
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Invalid status code: %d but expected 500", w.Code)
	}
}
