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
		{ID: "1", Username: "testuser1", Email: "testuser1@example.com"},
		{ID: "2", Username: "testuser2", Email: "testuser1@example.com"},
	}

	mockRepository.EXPECT().GetAllUsers().Return(want, nil)

	got := NewService(mockRepository).GetAllUsers()
	assert.ElementsMatch(t, want, got)
}

//
//func Test_GetUser(t *testing.T) {
//	id := "101"
//	want := "Get user 101"
//
//	got := NewSimpleService().GetUser(id)
//	assert.Equal(t, want, got)
//}
//
//func Test_CreateUser(t *testing.T) {
//	want := "Create new user"
//	got := NewSimpleService().CreateUser()
//	assert.Equal(t, want, got)
//}
//
//func Test_UpdateUser(t *testing.T) {
//	id := "101"
//	want := "Update user 101"
//	got := NewSimpleService().UpdateUser(id)
//	assert.Equal(t, want, got)
//}
//
//func Test_DeleteUser(t *testing.T) {
//	id := "101"
//	want := "Delete user 101"
//	got := NewSimpleService().DeleteUser(id)
//	assert.Equal(t, want, got)
//}
