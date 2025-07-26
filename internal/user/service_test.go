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
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockRepo := NewMockRepository(ctrl)
	service := NewService(mockRepo)
	return mockRepo, service
}

func TestService_CreateUser(t *testing.T) {
	tests := []struct {
		name       string
		req        *CreateUserRequest
		want       *model.User
		expectMock func(repo *MockRepository, want *model.User)
		wantErr    string
	}{
		{
			name: "success",
			req: &CreateUserRequest{
				Username: "testuser01",
				Email:    "testuser01@example.com",
			},
			want: model.NewUser("testuser01", "testuser01@example.com"),
			expectMock: func(repo *MockRepository, want *model.User) {
				repo.EXPECT().Create(gomock.Any()).Return(want, nil)
			},
			wantErr: "",
		},
		{
			name: "username has invalid format",
			req: &CreateUserRequest{
				Username: "01",
				Email:    "testuser01@example.com",
			},
			want:       nil,
			expectMock: nil,
			wantErr:    "invalid username: must be alphanumeric, at least 3 character",
		},
		{
			name: "db error",
			req: &CreateUserRequest{
				Username: "testuser01",
				Email:    "testuser01@example.com",
			},
			want: nil,
			expectMock: func(repo *MockRepository, want *model.User) {
				repo.EXPECT().Create(gomock.Any()).Return(nil, errors.New("db error"))
			},
			wantErr: "db error",
		},
		{
			name: "username taken",
			req: &CreateUserRequest{
				Username: "testuser01",
				Email:    "testuser01@example.com",
			},
			want: nil,
			expectMock: func(repo *MockRepository, want *model.User) {
				repo.EXPECT().Create(gomock.Any()).
					Return(nil, errors.New(`violates unique constraint "uni_users_username"`))
			},
			wantErr: "username already exists",
		},
		{
			name: "email taken",
			req: &CreateUserRequest{
				Username: "testuser01",
				Email:    "testuser01@example.com",
			},
			want: nil,
			expectMock: func(repo *MockRepository, want *model.User) {
				repo.EXPECT().Create(gomock.Any()).
					Return(nil, errors.New(`violates unique constraint "uni_users_email"`))
			},
			wantErr: "email already exists",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo, service := setupMockRepoAndService(t)

			if test.expectMock != nil {
				test.expectMock(mockRepo, test.want)
			}

			got, err := service.CreateUser(test.req)

			if test.wantErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			} else {
				assert.Nil(t, got)
				assert.ErrorContains(t, err, test.wantErr)
			}
		})
	}
}

func TestService_UpdateUser(t *testing.T) {
	id := uuid.MustParse("5faeb0b2-43b3-4a7a-aaf2-77b71ea59d90")
	tests := []struct {
		name       string
		req        *UpdateUserRequest
		old        *model.User
		want       *model.User
		expectMock func(repo *MockRepository, old, want *model.User)
		wantErr    string
	}{
		{
			name: "success",
			req: &UpdateUserRequest{
				Username: ptr("updatedtestuser01"),
				Email:    ptr("updatedtestuser01@example.com"),
			},
			old: &model.User{
				ID:       id,
				Username: "testuser01",
				Email:    "testuser01@example.com",
			},
			want: &model.User{
				ID:       id,
				Username: "updatedtestuser01",
				Email:    "updatedtestuser01@example.com",
			},
			expectMock: func(mockRepo *MockRepository, old, want *model.User) {
				mockRepo.EXPECT().FindByID(gomock.Any()).Return(old, nil)
				mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(nil, errors.New("user not found"))
				mockRepo.EXPECT().FindByEmail(gomock.Any()).Return(nil, errors.New("email not found"))
				mockRepo.EXPECT().Update(gomock.Any()).Return(want, nil)
			},
			wantErr: "",
		},
		{
			name: "user not found",
			req: &UpdateUserRequest{
				Username: ptr("updatedtestuser01"),
				Email:    ptr("updatedtestuser01@example.com"),
			},
			old: &model.User{
				ID:       uuid.New(),
				Username: "testuser01",
				Email:    "testuser01@example.com",
			},
			want: nil,
			expectMock: func(mockRepo *MockRepository, old, want *model.User) {
				mockRepo.EXPECT().FindByID(gomock.Any()).
					Return(nil, apperrors.NewNotFoundError("user", old.ID))
			},
			wantErr: apperrors.NewNotFoundError("user", id).Error(),
		},
		{
			name: "username already exists",
			req: &UpdateUserRequest{
				Username: ptr("updatedtestuser01"),
				Email:    ptr("updatedtestuser01@example.com"),
			},
			old: &model.User{
				ID:       uuid.New(),
				Username: "testuser01",
				Email:    "testuser01@example.com",
			},
			want: nil,
			expectMock: func(mockRepo *MockRepository, old, want *model.User) {
				mockRepo.EXPECT().FindByID(gomock.Any()).Return(old, nil)
				mockRepo.EXPECT().FindByUsername(gomock.Any()).
					Return(old, nil)
			},
			wantErr: "username already exists",
		},
		{
			name: "email already exists",
			req: &UpdateUserRequest{
				Username: ptr("updatedtestuser01"),
				Email:    ptr("updatedtestuser01@example.com"),
			},
			old: &model.User{
				ID:       uuid.New(),
				Username: "testuser01",
				Email:    "testuser01@example.com",
			},
			want: nil,
			expectMock: func(mockRepo *MockRepository, old, want *model.User) {
				mockRepo.EXPECT().FindByID(gomock.Any()).Return(old, nil)
				mockRepo.EXPECT().FindByUsername(gomock.Any()).Return(nil, errors.New("user not found"))
				mockRepo.EXPECT().FindByEmail(gomock.Any()).
					Return(&model.User{}, nil)
			},
			wantErr: "email already exists",
		},
		{
			name: "invalid username format",
			req: &UpdateUserRequest{
				Username: ptr("a1"),
				Email:    ptr("updatedtestuser01@example.com"),
			},
			old: &model.User{
				ID:       uuid.New(),
				Username: "testuser01",
				Email:    "testuser01@example.com",
			},
			want: nil,
			expectMock: func(mockRepo *MockRepository, old, want *model.User) {
				mockRepo.EXPECT().FindByID(gomock.Any()).Return(old, nil)
			},
			wantErr: "invalid username: must be alphanumeric, at least 3 character",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo, service := setupMockRepoAndService(t)

			if test.expectMock != nil {
				test.expectMock(mockRepo, test.old, test.want)
			}

			got, err := service.UpdateUser(id, test.req)

			if test.wantErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			} else {
				assert.Nil(t, got)
				assert.ErrorContains(t, err, test.wantErr)
			}
		})
	}
}

func TestService_GetUser(t *testing.T) {
	id := uuid.New()
	tests := []struct {
		name       string
		want       *model.User
		expectMock func(mockRepo *MockRepository, want *model.User)
		wantErr    string
	}{
		{
			name: "success",
			want: &model.User{
				ID:       id,
				Username: "testuser01",
				Email:    "testuser01@example.com",
			},
			expectMock: func(mockRepo *MockRepository, want *model.User) {
				mockRepo.EXPECT().FindByID(gomock.Any()).Return(want, nil)
			},
			wantErr: "",
		},
		{
			name: "user not found",
			want: nil,
			expectMock: func(mockRepo *MockRepository, want *model.User) {
				mockRepo.EXPECT().FindByID(gomock.Any()).
					Return(nil, errors.New("user not found"))
			},
			wantErr: apperrors.NewNotFoundError("user", id).Error(),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo, service := setupMockRepoAndService(t)
			if test.expectMock != nil {
				test.expectMock(mockRepo, test.want)
			}
			got, err := service.GetUser(id)
			if test.wantErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			} else {
				assert.Nil(t, got)
				assert.ErrorContains(t, err, test.wantErr)
			}

		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	id := uuid.New()
	tests := []struct {
		name       string
		id         uuid.UUID
		expectMock func(mockRepo *MockRepository)
		wantErr    string
	}{
		{
			name: "success",
			id:   id,
			expectMock: func(mockRepo *MockRepository) {
				mockRepo.EXPECT().Delete(gomock.Any()).Return(nil)
				mockRepo.EXPECT().FindByID(gomock.Any()).Return(&model.User{}, nil)
			},
			wantErr: "",
		},
		{
			name: "user not found",
			id:   id,
			expectMock: func(mockRepo *MockRepository) {
				mockRepo.EXPECT().FindByID(gomock.Any()).
					Return(&model.User{}, errors.New("user not found"))
			},
			wantErr: "not found",
		},
		{
			name: "deletion failed",
			id:   id,
			expectMock: func(mockRepo *MockRepository) {
				mockRepo.EXPECT().Delete(gomock.Any()).
					Return(errors.New("failed to delete user"))
				mockRepo.EXPECT().FindByID(gomock.Any()).
					Return(&model.User{}, nil)
			},
			wantErr: "failed to delete user",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo, service := setupMockRepoAndService(t)
			if test.expectMock != nil {
				test.expectMock(mockRepo)
			}

			err := service.DeleteUser(id)

			if test.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.wantErr)
			}
		})
	}
}

func TestService_GetUsers(t *testing.T) {
	tests := []struct {
		name       string
		want       []*model.User
		expectMock func(mockRepo *MockRepository, users []*model.User)
		wantErr    string
	}{
		{
			name: "success",
			want: []*model.User{},
			expectMock: func(mockRepo *MockRepository, users []*model.User) {
				mockRepo.EXPECT().FindAll().Return(users, nil)
			},
			wantErr: "",
		},
		{
			name: "failed",
			want: []*model.User{},
			expectMock: func(mockRepo *MockRepository, users []*model.User) {
				mockRepo.EXPECT().FindAll().Return(users, errors.New("failed"))
			},
			wantErr: "failed to get all users",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo, service := setupMockRepoAndService(t)

			if test.expectMock != nil {
				test.expectMock(mockRepo, test.want)
			}

			got, err := service.GetUsers()

			if test.wantErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, test.want, got)
			} else {
				assert.Nil(t, got)
				assert.ErrorContains(t, err, test.wantErr)
			}
		})
	}
}

func Test(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  string
	}{
		{
			name:     "success",
			username: "testuser01",
			wantErr:  "",
		},
		{
			name:     "is a number",
			username: "123",
			wantErr:  "invalid username: must not be a number",
		},
		{
			name:     "has less than 3 characters",
			username: "ab",
			wantErr:  "invalid username: must be alphanumeric, at least 3 character",
		},
		{
			name:     "has less than 2 letters",
			username: "a12",
			wantErr:  "invalid username: must have at least two letters",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := validateUsernameFormat(test.username)
			if test.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.ErrorContains(t, err, test.wantErr)
			}
		})
	}
}
