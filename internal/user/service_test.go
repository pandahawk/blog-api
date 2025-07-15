package user

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pandahawk/blog-api/internal/apperrors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ptr(s string) *string {
	return &s
}

func setupMockRepoAndService(t *testing.T) (*MockRepository, Service) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockRepo := NewMockRepository(ctrl)
	service := NewService(mockRepo)
	return mockRepo, service
}

//todo: replace error asserts with errorcontains

func TestService_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := CreateUserRequest{
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		wantUser := User{
			ID:       1001,
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(User{}, false)
		mockRepo.EXPECT().FindByEmail(gomock.Any()).Return(User{}, false)
		mockRepo.EXPECT().Save(gomock.Any()).Return(wantUser, nil)

		gotUser, err := service.CreateUser(req)

		assert.NoError(t, err)
		assert.Equal(t, wantUser, gotUser)
	})

	t.Run("username taken", func(t *testing.T) {
		req := CreateUserRequest{
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		wantUser := User{
			ID:       1001,
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(wantUser, true)

		_, err := service.CreateUser(req)

		assert.ErrorContains(t, err, "username already exists")
	})

	t.Run("email taken", func(t *testing.T) {
		req := CreateUserRequest{
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		wantUser := User{
			ID:       1001,
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(User{}, false)
		mockRepo.EXPECT().FindByEmail(gomock.Any()).Return(wantUser, true)

		_, err := service.CreateUser(req)

		assert.ErrorContains(t, err, "email already exists")
	})
}

func TestService_UpdateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		id := 1001
		req := UpdateUserRequest{
			Username: ptr("updatedtestuser01"),
			Email:    ptr("updatedtestuser01@example.com"),
		}
		oldUser := User{
			ID:       1001,
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		wantUser := User{
			ID:       1001,
			Username: "updatedtestuser01",
			Email:    "updatedtestuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().Update(gomock.Any()).Return(wantUser, nil)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(oldUser, true)

		gotUser, err := service.UpdateUser(id, req)

		assert.NoError(t, err)
		assert.Equal(t, wantUser, gotUser)
	})

	t.Run("username not found", func(t *testing.T) {
		id := 1001
		req := UpdateUserRequest{
			Username: ptr("updatedtestuser01"),
			Email:    ptr("updatedtestuser01@example.com"),
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(User{}, false)

		_, err := service.UpdateUser(id, req)

		assert.ErrorContains(t, err, fmt.Sprintf("user with ID %d not found", id))
	})

	t.Run("email not found", func(t *testing.T) {
		id := 1001
		req := UpdateUserRequest{
			Username: ptr("  "),
			Email:    ptr("updatedtestuser01@example.com"),
		}
		oldUser := User{
			ID:       1001,
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(oldUser, true)

		_, err := service.UpdateUser(id, req)

		assert.ErrorContains(t, err, "username can not be empty")
	})

	t.Run("email empty", func(t *testing.T) {
		id := 1001
		req := UpdateUserRequest{
			Username: ptr("updatedtestuser01"),
			Email:    ptr(" "),
		}
		oldUser := User{
			ID:       1001,
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(oldUser, true)

		_, err := service.UpdateUser(id, req)

		assert.ErrorContains(t, err, "email cannot be empty")
	})
}

func TestService_GetUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		id := 1001
		wantUser := User{
			ID:       1001,
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(wantUser, true)

		gotUser, err := service.GetUser(id)

		assert.NoError(t, err)
		assert.Equal(t, wantUser, gotUser)
	})

	t.Run("user not found", func(t *testing.T) {
		id := 1001
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(User{}, false)

		_, err := service.GetUser(id)

		assert.ErrorContains(t, err, apperrors.NewNotFoundError("user", id).Error())
	})

}

func TestService_DeleteUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		id := 1001
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().Delete(gomock.Any()).Return(true)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(User{}, true)

		err := service.DeleteUser(id)

		assert.NoError(t, err)
	})
	t.Run("user not found", func(t *testing.T) {
		id := 1001
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(User{}, false)

		err := service.DeleteUser(id)

		assert.ErrorContains(t, err, apperrors.NewNotFoundError("user", id).Error())
	})
	t.Run("deletion failed", func(t *testing.T) {
		id := 1001
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().Delete(gomock.Any()).Return(false)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(User{}, true)

		err := service.DeleteUser(id)

		//todo: check why this error string is fine although err has a different message
		assert.ErrorContains(t, err, "failed to delete user")
	})
}

func TestService_GetAllUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindAll().Return([]User{}, true)
	})
}
