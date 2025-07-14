package user

//func Test_GetAllUsers(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//	mockRepository := NewMockRepository(ctrl)
//
//	want := []User{
//		{ID: 1, Username: "testuser1", Email: "testuser1@example.com"},
//		{ID: 2, Username: "testuser2", Email: "testuser1@example.com"},
//	}
//
//	mockRepository.EXPECT().FindAll().Return(want, nil)
//
//	got := NewService(mockRepository).GetAllUsers()
//	assert.ElementsMatch(t, want, got)
//}
