package internal

import "fmt"

type User struct {
	uid      uint32
	free     bool
	username string
}

type SystemUsers struct {
	users []*User
}

func NewSystemUser(amount uint32) *SystemUsers {
	users := make([]*User, amount)

	for i := range users {
		users[i] = &User{
			uid:      uint32(i),
			free:     true,
			username: fmt.Sprintf("user%d", i),
		}
	}

	return &SystemUsers{
		users,
	}
}

func (su *SystemUsers) Acquire() (*User, error) {
	for _, user := range su.users {
		if user.free {
			user.free = false
			return user, nil
		}
	}

	return &User{}, fmt.Errorf("free user not found")
}

func (su *SystemUsers) Release(uid uint32) {
	for _, user := range su.users {
		if user.uid == uid {
			user.free = true
		}
	}
}
