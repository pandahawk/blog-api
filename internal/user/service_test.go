package user

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func ptr(s string) *string {
	return &s
}

func setupMockRepo(t *testing.T) *MockRepository {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	return NewMockRepository(ctrl)
}

func TestService_CreateUser(t *testing.T) {
	mockRepo := setupMockRepo(t)
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
		mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(User{}, false)
		mockRepo.EXPECT().FindByEmail(gomock.Any()).Return(User{}, false)
		mockRepo.EXPECT().Save(gomock.Any()).Return(wantUser, nil)

		service := NewService(mockRepo)

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
		mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(wantUser, true)

		service := NewService(mockRepo)

		_, err := service.CreateUser(req)
		assert.Errorf(t, err, "username already exists")
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
		mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(User{}, false)
		mockRepo.EXPECT().FindByEmail(gomock.Any()).Return(wantUser, true)

		service := NewService(mockRepo)

		_, err := service.CreateUser(req)
		assert.Errorf(t, err, "email already exists")
	})
}

func TestService_UpdateUser(t *testing.T) {
	mockRepo := setupMockRepo(t)

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
		mockRepo.EXPECT().Update(gomock.Any()).Return(wantUser, nil)
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(oldUser, true)
		service := NewService(mockRepo)

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
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(User{}, false)
		service := NewService(mockRepo)

		_, err := service.UpdateUser(id, req)

		assert.Errorf(t, err, fmt.Sprintf("user with id %d not found", id))
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
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(oldUser, true)
		service := NewService(mockRepo)

		_, err := service.UpdateUser(id, req)

		assert.Errorf(t, err, "username cannot be empty")
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
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(oldUser, true)
		service := NewService(mockRepo)

		_, err := service.UpdateUser(id, req)

		assert.Errorf(t, err, "email cannot be empty")
	})
}

func TestService_GetUser(t *testing.T) {
	mockRepo := setupMockRepo(t)
	t.Run("success", func(t *testing.T) {
		id := 1001
		wantUser := User{
			ID:       1001,
			Username: "testuser01",
			Email:    "testuser01@example.com",
		}
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(wantUser, true)
		service := NewService(mockRepo)

		gotUser, err := service.GetUser(id)

		assert.NoError(t, err)
		assert.Equal(t, wantUser, gotUser)
	})

	t.Run("user not found", func(t *testing.T) {
		id := 1001
		mockRepo.EXPECT().FindByID(gomock.Any()).Return(User{}, false)
		service := NewService(mockRepo)

		_, err := service.GetUser(id)

		assert.Errorf(t, err, fmt.Sprintf("user %d not found", id))
	})

}
