package main

import (
	"github.com/32leaves/bel"
)

type User struct {
	Name string
}
type AddUserRequest struct {
	NewUser User
}
type UserService interface {
	AddUser(AddUserRequest) error
}

func FollowStructs() {
	ts, err := bel.Extract((*UserService)(nil), bel.FollowStructs)
	if err != nil {
		panic(err)
	}

	err = bel.Render(ts)
	if err != nil {
		panic(err)
	}
}
