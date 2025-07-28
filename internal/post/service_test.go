package post

import (
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/pandahawk/blog-api/internal/apperrors"
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

func TestService_CreatePost(t *testing.T) {
	tests := []struct {
		name          string
		req           *CreatePostRequest
		want          *model.Post
		mockBehaviour func(repo *MockRepository, post *model.Post)
		wantErr       string
	}{
		{
			name: "success",
			req: &CreatePostRequest{
				Title:    "test title",
				Content:  "test content",
				AuthorID: uuid.UUID{},
			},
			want: &model.Post{
				Title:   "test title",
				Content: "test content",
				UserID:  uuid.UUID{},
				User: &model.User{
					ID:       uuid.UUID{},
					Username: "testuser",
					Email:    "testuser@mail.com",
				},
			},
			mockBehaviour: func(repo *MockRepository, post *model.Post) {
				repo.EXPECT().Create(gomock.Any()).Return(post, nil)
			},
			wantErr: "",
		},
		{
			name: "author id not found",
			req: &CreatePostRequest{
				Title:    "test title",
				Content:  "test content",
				AuthorID: uuid.UUID{},
			},
			want: nil,
			mockBehaviour: func(repo *MockRepository, post *model.Post) {
				repo.EXPECT().Create(gomock.Any()).Return(nil, errors.New("author not found"))
			},
			wantErr: "author not found",
		},
		{
			name: "blank tile",
			req: &CreatePostRequest{
				Title:    " ",
				Content:  "test content",
				AuthorID: uuid.UUID{},
			},
			want:          nil,
			mockBehaviour: nil,
			wantErr:       "title must not be blank",
		},
		{
			name: "blank content",
			req: &CreatePostRequest{
				Title:    "test title",
				Content:  " ",
				AuthorID: uuid.UUID{},
			},
			want:          nil,
			mockBehaviour: nil,
			wantErr:       "content must not be blank",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo, service := setup(t)
			if test.mockBehaviour != nil {
				test.mockBehaviour(mockRepo, test.want)
			}
			got, err := service.CreatePost(test.req)
			if test.wantErr == "" {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, test.want.Title, got.Title)
				assert.Equal(t, test.want.Content, got.Content)
			} else {
				assert.Nil(t, got)
				assert.ErrorContains(t, err, test.wantErr)
			}
		})
	}
}

func TestService_GetPost(t *testing.T) {
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
			searchID: uuid.Nil,
			wantPost: nil,
			expectMock: func(mockRepo *MockRepository, wantPost *model.Post) {
				mockRepo.EXPECT().FindByID(gomock.Any()).
					Return(wantPost, apperrors.NewNotFoundError("post", uuid.Nil))
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

func TestService_GetPosts(t *testing.T) {
	tests := []struct {
		name          string
		want          []*model.Post
		mockBehaviour func(repo *MockRepository, posts []*model.Post)
	}{
		// TODO: test cases
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
