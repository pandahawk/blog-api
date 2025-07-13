package user

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetAllUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepository := NewMockRepository(ctrl)

	want := []User{
		{ID: 1, Username: "testuser1", Email: "testuser1@example.com"},
		{ID: 2, Username: "testuser2", Email: "testuser1@example.com"},
	}

	mockRepository.EXPECT().GetAllUsers().Return(want, nil)

	got := NewService(mockRepository).GetAllUsers()
	assert.ElementsMatch(t, want, got)
}
