package memory

import (
	"sync"

	"github.com/ltick/dummy/app/domain/user"
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

func (r *userRepository) FindAll() ([]*user.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	users := make([]*user.User, len(r.users))
	i := 0
	for _, user := range r.users {
		users[i] = user.NewUser(user.ID, user.Email)
		i++
	}
	return users, nil
}

func (r *userRepository) FindByEmail(email string) (*user.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, user := range r.users {
		if user.Email == email {
			return user.NewUser(user.ID, user.Email), nil
		}
	}
	return nil, nil
}

func (r *userRepository) Save(user *user.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.GetID()] = &User{
		ID:    user.GetID(),
		Email: user.GetEmail(),
	}
	return nil
}

type User struct {
	ID    string
	Email string
}