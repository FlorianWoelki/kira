package pool

import "fmt"

type User struct {
	Uid      uint32
	Free     bool
	Username string
}

type SystemUsers struct {
	users []*User
}

func (su SystemUsers) GetUsers() []*User {
	return su.users
}

func NewSystemUser(amount uint32) *SystemUsers {
	users := make([]*User, amount)

	for i := range users {
		users[i] = &User{
			Uid:      uint32(i),
			Free:     true,
			Username: fmt.Sprintf("user%d", i),
		}
	}

	return &SystemUsers{
		users,
	}
}

func (su *SystemUsers) Acquire() (*User, error) {
	for _, user := range su.users {
		if user.Free {
			user.Free = false
			return user, nil
		}
	}

	return &User{}, fmt.Errorf("free user not found")
}

func (su *SystemUsers) Release(uid uint32) {
	for _, user := range su.users {
		if user.Uid == uid {
			user.Free = true
		}
	}
}
