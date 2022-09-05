package mock1

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
)

func TestCreateRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repo := NewMockRepository(ctrl)
	{
		repo.EXPECT().Create(gomock.Eq("name"), gomock.Eq([]byte("fly"))).Return(nil)
		repo.EXPECT().Retrieve(gomock.Eq("name")).Return([]byte("fly"), nil)
		val, err := CreateRepo(repo, "name", "fly")
		if val != "fly" || err != nil {
			t.Fatal("CreateRepo failed")
		}
	}
	{
		repo.EXPECT().Create(gomock.Eq("name"), gomock.Eq([]byte("cc"))).Return(errors.New("Create err"))
		val, err := CreateRepo(repo, "name", "cc")
		if err != nil && val == "" {
			t.Log("Create err test success!")
		} else {
			t.Fatal("Create err test failed!")
		}
	}
	{
		repo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil).Times(5)
		repo.EXPECT().Retrieve(gomock.Any()).Return([]byte(""), nil).Times(5)
		val, err := CreateRepo(repo, "zzz", "bb")
		if err != nil && val == "" {
			t.Log("Create err test success!")
		} else {
			t.Fatal("Create err test failed!")
		}
	}
}
