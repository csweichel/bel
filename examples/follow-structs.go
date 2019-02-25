package main

import (
	"github.com/32leaves/bel"
)

// User is a struct describing users
type User struct {
	Name string
}

// AddUserRequest is the single prameter to create users
type AddUserRequest struct {
	NewUser User
}

// UserService enables the creation of users
type UserService interface {
	AddUser(AddUserRequest) error
}

// FollowStructs demonstrates the use of the bel.FollowStructs config option
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

func init() {
	examples["follow-structs"] = FollowStructs
}
