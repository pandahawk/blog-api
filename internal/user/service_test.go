package user

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/apperrors"
	"github.com/pandahawk/blog-api/internal/shared/model"
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
		wantUser := model.NewUser("testuser01", "testuser01@example.com")
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().Create(gomock.Any()).Return(wantUser, nil)

		gotUser, err := service.CreateUser(&req)

		assert.NoError(t, err)
		assert.Equal(t, wantUser, gotUser)
	})
	t.Run("username has invalid format", func(t *testing.T) {
		req := CreateUserRequest{
			Username: "01",
			Email:    "testuser01@example.com",
		}
		_, service := setupMockRepoAndService(t)

		_, err := service.CreateUser(&req)

		assert.ErrorContains(t, err, "invalid username: must be alphanumeric, at least 3 character")
	})

	t.Run("db error", func(t *testing.T) {
		req := CreateUserRequest{
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().Create(gomock.Any()).Return(nil, errors.New(""))

		_, err := service.CreateUser(&req)

		assert.Error(t, err)
	})

	t.Run("username taken", func(t *testing.T) {
		req := CreateUserRequest{
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().Create(gomock.Any()).
			Return(nil, errors.New(`violates unique constraint "uni_users_username"`))
		_, err := service.CreateUser(&req)

		assert.ErrorContains(t, err, "username already exists")
	})

	t.Run("email taken", func(t *testing.T) {
		req := CreateUserRequest{
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().Create(gomock.Any()).
			Return(nil, errors.New(`violates unique constraint "uni_users_email"`))

		_, err := service.CreateUser(&req)

		assert.ErrorContains(t, err, "email already exists")
	})
}

func TestService_UpdateUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		req := UpdateUserRequest{
			Username: ptr("updatedtestuser01"),
			Email:    ptr("updatedtestuser01@example.com"),
		}
		oldUser := model.User{
			ID:       uuid.New(),
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		wantUser := model.User{
			ID:       oldUser.ID,
			Username: "updatedtestuser01",
			Email:    "updatedtestuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(&oldUser, nil)
		mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(nil, errors.New("user not found"))
		mockRepo.EXPECT().FindByEmail(gomock.Any()).Return(nil, errors.New("email not found"))
		mockRepo.EXPECT().Update(gomock.Any()).Return(&wantUser, nil)
		gotUser, err := service.UpdateUser(oldUser.ID, &req)

		assert.NoError(t, err)
		assert.Equal(t, &wantUser, gotUser)
	})

	t.Run("user not found", func(t *testing.T) {
		req := UpdateUserRequest{
			Username: ptr("updatedtestuser01"),
			Email:    ptr("updatedtestuser01@example.com"),
		}
		oldUser := model.User{
			ID:       uuid.New(),
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).
			Return(nil, apperrors.NewNotFoundError("user", oldUser.ID))

		_, err := service.UpdateUser(oldUser.ID, &req)

		assert.ErrorContains(t, err, "not found")
	})

	t.Run("username already exists", func(t *testing.T) {
		oldUser := model.User{
			ID:       uuid.New(),
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		req := UpdateUserRequest{
			Username: ptr("updatedtestuser01"),
			Email:    ptr("updatedtestuser01@example.com"),
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(&oldUser, nil)
		mockRepo.EXPECT().FindByUsername(gomock.Any()).
			Return(&oldUser, nil)
		_, err := service.UpdateUser(oldUser.ID, &req)

		assert.ErrorContains(t, err, "username already exists")
	})
	t.Run("email already exists", func(t *testing.T) {
		oldUser := model.User{
			ID:       uuid.New(),
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		req := UpdateUserRequest{
			Username: ptr("updatedtestuser01"),
			Email:    ptr("updatedtestuser01@example.com"),
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(&oldUser, nil)
		mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(nil, errors.New("user not found"))
		mockRepo.EXPECT().FindByEmail(gomock.Any()).
			Return(&model.User{}, nil)
		_, err := service.UpdateUser(oldUser.ID, &req)

		assert.ErrorContains(t, err, "email already exists")
	})

	t.Run("invalid username format", func(t *testing.T) {
		req := UpdateUserRequest{
			Username: ptr("a1"),
			Email:    ptr("updatedtestuser01@example.com"),
		}
		oldUser := model.User{
			ID:       uuid.New(),
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(&oldUser, nil)

		_, err := service.UpdateUser(oldUser.ID, &req)

		assert.ErrorContains(t, err, "invalid username: must be alphanumeric, at least 3 character")
	})
}

func TestService_GetUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		wantUser := model.User{
			ID:       uuid.New(),
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(&wantUser, nil)

		gotUser, err := service.GetUser(wantUser.ID)

		assert.NoError(t, err)
		assert.Equal(t, &wantUser, gotUser)
	})

	t.Run("user not found", func(t *testing.T) {
		id := uuid.New()
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).
			Return(nil, errors.New("user not found"))

		_, err := service.GetUser(id)

		assert.ErrorContains(t, err,
			apperrors.NewNotFoundError("user", id).Error())
	})

}

func TestService_DeleteUser(t *testing.T) {

	t.Run("success", func(t *testing.T) {

		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().Delete(gomock.Any()).Return(nil)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(&model.User{}, nil)

		err := service.DeleteUser(uuid.New())

		assert.NoError(t, err)
	})
	t.Run("user not found", func(t *testing.T) {
		id := uuid.New()
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(&model.User{}, errors.New("user not found"))

		err := service.DeleteUser(id)

		assert.ErrorContains(t, err, apperrors.NewNotFoundError("user", id).Error())
	})
	t.Run("deletion failed", func(t *testing.T) {
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().Delete(gomock.Any()).Return(errors.New("user not found"))
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(&model.User{}, nil)

		err := service.DeleteUser(uuid.New())

		assert.ErrorContains(t, err, "failed to delete user")
	})
}

func TestService_GetAllUsers(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var wantUsers []*model.User
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindAll().Return(wantUsers, nil)

		gotUsers, err := service.GetUsers()

		assert.NoError(t, err)
		assert.Equal(t, wantUsers, gotUsers)
	})
	t.Run("failed", func(t *testing.T) {
		var wantUsers []*model.User
		mockRepo, service := setupMockRepoAndService(t)
		mockRepo.EXPECT().FindAll().Return(wantUsers, errors.New("failed"))

		_, err := service.GetUsers()

		assert.ErrorContains(t, err, "failed to get all users")
	})
}

func TestValidateUsernameFormat(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		err := validateUsernameFormat("testuser01")
		assert.NoError(t, err)
	})
	t.Run("is a number", func(t *testing.T) {
		err := validateUsernameFormat("123")
		var ie *apperrors.InvalidInputError
		assert.True(t, errors.As(err, &ie))
	})
	t.Run("has less than 3 characters", func(t *testing.T) {
		err := validateUsernameFormat("ab")
		var ie *apperrors.InvalidInputError
		assert.True(t, errors.As(err, &ie))
	})
	t.Run("has less than 3 characters", func(t *testing.T) {
		err := validateUsernameFormat("ab")
		var ie *apperrors.InvalidInputError
		assert.True(t, errors.As(err, &ie))
	})
	t.Run("has less than 2 letters", func(t *testing.T) {
		err := validateUsernameFormat("a12")
		var ie *apperrors.InvalidInputError
		assert.True(t, errors.As(err, &ie))
	})
}
