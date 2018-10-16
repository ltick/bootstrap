package http

import (
	"github.com/ltick/tick-framework"
	"github.com/ltick/dummy/app"
)

type userService struct {
	user app.User
}

func NewUserService(user app.User) *userService {
	return &userService{
		user: user,
	}
}

func (s *userService) ListUser(ctx *ltick.Context) error {
	users, err := s.user.ListUser()
	if err != nil {
		return err
	}

	res := &ListUserResponseType{
		Users: toUsers(users),
	}

	ctx.Write(res.Users)

	return nil
}

func (s *userService) RegisterUser(ctx *ltick.Context) error {
	if err := s.user.RegisterUser(ctx.Query("Email")); err != nil {
		return err
	}
	return nil
}

func toUsers(users []*app.User) []*User {
	res := make([]*User, len(users))
	for i, user := range users {
		res[i] = &User{
			Id:    user.ID,
			Email: user.Email,
		}
	}
	return res
}