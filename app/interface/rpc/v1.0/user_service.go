package v1

import (
	"context"

	"github.com/ltick/dummy/app/interface/rpc/v1.0/protocol"
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

func (s *userService) ListUser(ctx context.Context, in *protocol.ListUserRequestType) (*protocol.ListUserResponseType, error) {
	users, err := s.user.ListUser()
	if err != nil {
		return nil, err
	}

	res := &protocol.ListUserResponseType{
		Users: toUser(users),
	}

	return res, nil
}

func (s *userService) RegisterUser(ctx context.Context, in *protocol.RegisterUserRequestType) (*protocol.RegisterUserResponseType, error) {
	if err := s.user.RegisterUser(in.GetEmail()); err != nil {
		return &protocol.RegisterUserResponseType{}, err
	}
	return &protocol.RegisterUserResponseType{}, nil
}

func toUser(users []*app.User) []*protocol.User {
	res := make([]*protocol.User, len(users))
	for i, user := range users {
		res[i] = &protocol.User{
			Id:    user.ID,
			Email: user.Email,
		}
	}
	return res
}