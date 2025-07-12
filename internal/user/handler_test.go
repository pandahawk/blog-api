package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestRouter(service UserService) *gin.Engine {
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
	mockService := NewMockUserService(ctrl)
	mockService.EXPECT().
		GetAllUsers().
		Return("Get all users")

	router := setupTestRouter(mockService)

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "Get all users"}`, w.Body.String())
}

func Test_getUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockUserService(ctrl)
	mockService.EXPECT().
		GetUser(gomock.Any()).
		Return("Get user 101")
	router := setupTestRouter(mockService)

	id := "101"
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", id), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	expected := fmt.Sprintf(`{"message": "Get user %s"}`, id)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
}

func Test_createUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockUserService(ctrl)
	mockService.EXPECT().
		CreateUser().Return("Create new user")

	router := setupTestRouter(mockService)
	req, _ := http.NewRequest(http.MethodPost, "/users", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Create new user")
}

func Test_UpdateUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockUserService(ctrl)
	mockService.EXPECT().
		UpdateUser(gomock.Any()).Return("Update user 101")

	router := setupTestRouter(mockService)
	id := "101"
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%s", id), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Update user")
	assert.Contains(t, w.Body.String(), id)
}

func Test_deleteUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := NewMockUserService(ctrl)
	mockService.EXPECT().
		DeleteUser(gomock.Any()).Return("Delete user 101")
	router := setupTestRouter(mockService)
	id := "101"
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", id), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	expected := `{"message": "Delete user 101"}`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())

}
