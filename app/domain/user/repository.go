package user

type UserRepository interface {
	FindAll() ([]*User, error)
	FindByEmail(email string) (*User, error)
	Save(*User) error
}