package user

import (
	"testing"
	"sync"

	"github.com/stretchr/testify/assert"
)

type userRepository struct {
	mu    *sync.Mutex
	users map[string]*User
}

func NewUserRepository() *userRepository {
	return &userRepository{
		mu:    &sync.Mutex{},
		users: map[string]*User{},
	}
}

func (r *userRepository) FindAll() ([]*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	users := make([]*User, len(r.users))
	i := 0
	for _, u := range r.users {
		users[i] = NewUser(u.id, u.email)
		i++
	}
	return users, nil
}

func (r *userRepository) FindByemail(email string) (*User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.email == email {
			return NewUser(u.id, u.email), nil
		}
	}
	return nil, nil
}

func (r *userRepository) Save(user *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.GetID()] = &User{
		id:    user.GetID(),
		email: user.GetEmail(),
	}
	return nil
}

func TestNewService(t *testing.T) {
	repo := NewUserRepository()
	u := &User{
		id:    "1",
		email: "sam42@outlook.com",
	}
	err := repo.Save(u)
	assert.Nil(t, err)
	users, err := repo.FindAll()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(users))
	user, err := repo.FindByemail("sam42@outlook.com")
	assert.Nil(t, err)
	assert.Equal(t, "1", user.GetID())
	assert.Equal(t, "sam42@outlook.com", user.GetEmail())
}