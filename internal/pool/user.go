package pool

import "fmt"

type User struct {
	Uid uint32
	// Free indicates a flag whether the user is free or not.
	Free     bool
	Username string
}

// SystemUsers represents a pool of users.
type SystemUsers struct {
	users []*User
}

// NewSystemUser creates a new pool of users with a specific amount. By default all users
// are free.
func NewSystemUsers(amount uint32) *SystemUsers {
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

// GetUsers returns the list of users in the current pool.
func (su SystemUsers) GetUsers() []*User {
	return su.users
}

// Acquire returns a free user and marks it as used. If no user was found, it will return
// an error.
func (su *SystemUsers) Acquire() (*User, error) {
	for _, user := range su.users {
		if user.Free {
			user.Free = false
			return user, nil
		}
	}

	return &User{}, fmt.Errorf("free user not found (users: %v+)", su.users)
}

// Release marks a user with the specified `uid“ as free.
func (su *SystemUsers) Release(uid uint32) {
	for _, user := range su.users {
		if user.Uid == uid {
			user.Free = true
		}
	}
}
