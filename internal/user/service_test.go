package user

import (
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
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

func TestService_CreateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := CreateUserRequest{
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		wantUser := User{
			ID:       uuid.New(),
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(User{}, errors.New(""))
		mockRepo.EXPECT().FindByEmail(gomock.Any()).Return(User{}, errors.New(""))
		mockRepo.EXPECT().Create(gomock.Any()).Return(wantUser, nil)

		gotUser, err := service.CreateUser(req)

		assert.NoError(t, err)
		assert.Equal(t, wantUser, gotUser)
	})
	t.Run("db error", func(t *testing.T) {
		req := CreateUserRequest{
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(User{}, errors.New(""))
		mockRepo.EXPECT().FindByEmail(gomock.Any()).Return(User{}, errors.New(""))
		mockRepo.EXPECT().Create(gomock.Any()).Return(User{}, errors.New(""))

		_, err := service.CreateUser(req)

		assert.Error(t, err)
	})

	t.Run("username taken", func(t *testing.T) {
		req := CreateUserRequest{
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		wantUser := User{
			ID:       uuid.New(),
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(wantUser, nil)

		_, err := service.CreateUser(req)

		assert.ErrorContains(t, err, "username already exists")
	})

	t.Run("email taken", func(t *testing.T) {
		req := CreateUserRequest{
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		wantUser := User{
			ID:       uuid.New(),
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(User{}, errors.New(""))
		mockRepo.EXPECT().FindByEmail(gomock.Any()).Return(wantUser, nil)

		_, err := service.CreateUser(req)

		assert.ErrorContains(t, err, "email already exists")
	})
}

func TestService_UpdateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := UpdateUserRequest{
			Username: ptr("updatedtestuser01"),
			Email:    ptr("updatedtestuser01@example.com"),
		}
		oldUser := User{
			ID:       uuid.New(),
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		wantUser := User{
			ID:       oldUser.ID,
			Username: "updatedtestuser01",
			Email:    "updatedtestuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().Update(gomock.Any()).Return(wantUser, nil)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(oldUser, nil)

		gotUser, err := service.UpdateUser(oldUser.ID, req)

		assert.NoError(t, err)
		assert.Equal(t, wantUser, gotUser)
	})

	t.Run("user not found", func(t *testing.T) {
		id := uuid.New()
		req := UpdateUserRequest{
			Username: ptr("updatedtestuser01"),
			Email:    ptr("updatedtestuser01@example.com"),
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(User{}, errors.New("user not found"))

		_, err := service.UpdateUser(id, req)

		assert.ErrorContains(t, err, fmt.Sprintf("user with ID %d not found", id))
	})

	t.Run("email not found", func(t *testing.T) {
		req := UpdateUserRequest{
			Username: ptr("  "),
			Email:    ptr("updatedtestuser01@example.com"),
		}
		oldUser := User{
			ID:       uuid.New(),
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(oldUser, nil)

		_, err := service.UpdateUser(oldUser.ID, req)

		assert.ErrorContains(t, err, "username can not be empty")
	})

	t.Run("email empty", func(t *testing.T) {
		req := UpdateUserRequest{
			Username: ptr("updatedtestuser01"),
			Email:    ptr(" "),
		}
		oldUser := User{
			ID:       uuid.New(),
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(oldUser, nil)

		_, err := service.UpdateUser(oldUser.ID, req)

		assert.ErrorContains(t, err, "email can not be empty")
	})
}

func TestService_GetUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		wantUser := User{
			ID:       uuid.New(),
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(wantUser, nil)

		gotUser, err := service.GetUser(wantUser.ID)

		assert.NoError(t, err)
		assert.Equal(t, wantUser, gotUser)
	})

	t.Run("user not found", func(t *testing.T) {
		id := uuid.New()
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(User{}, errors.New("user not found"))

		_, err := service.GetUser(id)

		assert.ErrorContains(t, err, apperrors.NewNotFoundError("user", id).Error())
	})

}

func TestService_DeleteUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().Delete(gomock.Any()).Return(nil)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(User{}, nil)

		err := service.DeleteUser(uuid.New())

		assert.NoError(t, err)
	})
	t.Run("user not found", func(t *testing.T) {
		id := uuid.New()
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(User{}, errors.New("user not found"))

		err := service.DeleteUser(id)

		assert.ErrorContains(t, err, apperrors.NewNotFoundError("user", id).Error())
	})
	t.Run("deletion failed", func(t *testing.T) {
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().Delete(gomock.Any()).Return(errors.New("user not found"))
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(User{}, nil)

		err := service.DeleteUser(uuid.New())

		assert.ErrorContains(t, err, "failed to delete user")
	})
}

func TestService_GetAllUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var wantUsers []User
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindAll().Return(wantUsers, nil)

		gotUsers, err := service.GetAllUsers()

		assert.NoError(t, err)
		assert.Equal(t, wantUsers, gotUsers)
	})
	t.Run("failed", func(t *testing.T) {
		var wantUsers []User
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindAll().Return(wantUsers, errors.New("failed"))

		_, err := service.GetAllUsers()

		assert.ErrorContains(t, err, "failed to get all users")
	})
}
