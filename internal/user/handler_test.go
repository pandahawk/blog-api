package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_getAllUsers(t *testing.T) {

	gin.SetMode(gin.TestMode)

	router := gin.Default()
	userGroup := router.Group("/users")
	RegisterRoutes(userGroup)

	req, _ := http.NewRequest(http.MethodGet, "/users", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"message": "get all users"}`, w.Body.String())
}

func Test_getUserById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	userGroup := router.Group("/users")
	RegisterRoutes(userGroup)

	id := "101"
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/users/%s", id), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	expected := fmt.Sprintf(`{"message": "get user %s"}`, id)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())
}

func Test_createUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	userGroup := router.Group("/users")
	RegisterRoutes(userGroup)
	req, _ := http.NewRequest(http.MethodPost, "/users", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "create new user")
}

func Test_UpdateUserById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	userGroup := router.Group("/users")
	RegisterRoutes(userGroup)

	id := "101"
	req, _ := http.NewRequest(http.MethodPatch, fmt.Sprintf("/users/%s", id), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "update user")
	assert.Contains(t, w.Body.String(), id)
}

func Test_deleteUserById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	userGroup := router.Group("/users")
	RegisterRoutes(userGroup)

	id := "101"
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("/users/%s", id), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	expected := fmt.Sprintf(`{"message": "delete user %s"}`, id)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, w.Body.String())

}
