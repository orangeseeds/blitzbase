package store

import "testing"

func TestUserModel(t *testing.T) {

	s := NewStorage("../test.db", "../migrations")
	s.Connect()

	u := User{
		Username: "new_user",
		Email:    "new@mail.com",
		Password: "testing123",
	}

	if err := s.DB.DB().Ping(); err != nil {
		t.Error(err.Error())
	}
	if err := s.DB.Model(&u).Insert(); err != nil {
		t.Error(err.Error())
	}

}
