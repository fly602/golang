package gotest

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestGetFromDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDB(ctrl)
	m.EXPECT().Get(gomock.Eq("Tom")).Return(-2, errors.New("not exist"))
	m.EXPECT().Get(gomock.Any()).Return(630, nil).Times(2)
	m.EXPECT().Get(gomock.Not("Sam")).Return(0, nil)
	m.EXPECT().Get(gomock.Nil().String()).Return(0, errors.New("nil"))

	if v := GetFromDB(m, "Tom"); v != -1 {
		t.Fatal("expected -1, but got", v)
	}
	if v := GetFromDB(m, "CC"); v != 630 {
		t.Fatal("expected -1, but got", v)
	}
	if v := GetFromDB(m, "Sam"); v != 630 {
		t.Fatal("expected -1, but got", v)
	}
	if v := GetFromDB(m, "DD"); v != 0 {
		t.Fatal("expected -1, but got", v)
	}
	if v := GetFromDB(m, gomock.Nil().String()); v != -1 {
		t.Fatal("expected -1, but got", v)
	}
}
