package post

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/shared/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup(t *testing.T) (*MockRepository, Service) {
	t.Helper()
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)
	mockRepo := NewMockRepository(ctrl)
	service := NewService(mockRepo)
	return mockRepo, service
}

//func TestPostService_CreatePost(t *testing.T) {
//	t.Run("success", func(t *testing.T) {
//		mockRepo, service := setup(t)
//		wantUser := &model.User{
//			ID:       uuid.UUID{},
//			Username: "testuser",
//			Email:    "testuser@mail.test",
//		}
//		wantPost := &model.Post{
//			ID:      uuid.New(),
//			Title:   "First Post",
//			Content: "This is a test gotPost",
//			UserID:  wantUser.ID,
//			User:    wantUser,
//		}
//		req := &CreatePostRequest{
//			Title:    wantPost.Title,
//			Content:  wantPost.Content,
//			AuthorID: wantPost.UserID,
//		}
//		mockRepo.EXPECT().Create(gomock.Any()).Return(wantPost, nil)
//
//		gotPost, err := service.CreatePost(req)
//
//		assert.NoError(t, err)
//		assert.Equal(t, wantPost, gotPost)
//	})
//}

func Test_service_GetPost(t *testing.T) {
	tests := []struct {
		name       string
		searchID   uuid.UUID
		expectMock func(mockRepo *MockRepository, wantPost *model.Post)
		wantPost   *model.Post
		wantErr    string
	}{
		{
			name:     "success",
			searchID: uuid.New(),
			wantPost: &model.Post{
				ID:      uuid.New(),
				Title:   "First Post",
				Content: "This is a test gotPost",
			},
			expectMock: func(mockRepo *MockRepository, wantPost *model.Post) {
				mockRepo.EXPECT().FindByID(gomock.Any()).Return(wantPost, nil)
			},
			wantErr: "",
		},
		{
			name:     "post not found",
			searchID: uuid.New(),
			wantPost: nil,
			expectMock: func(mockRepo *MockRepository, wantPost *model.Post) {
				mockRepo.EXPECT().FindByID(gomock.Any()).
					Return(wantPost, fmt.Errorf("not found"))
			},
			wantErr: "not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo, service := setup(t)

			if tt.expectMock != nil {
				tt.expectMock(mockRepo, tt.wantPost)
			}
			gotPost, err := service.GetPost(tt.searchID)
			if tt.wantErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantPost, gotPost)
			} else {
				assert.Nil(t, gotPost)
				assert.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}
