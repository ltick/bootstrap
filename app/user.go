package app

import (
	"context"
	"fmt"

	"github.com/ltick/dummy/app/domain/user"
	"github.com/ltick/tick-framework/config"
	"github.com/ltick/tick-framework/utility"
	"github.com/satori/go.uuid"
)

type UserInterface interface {
	ListUser() ([]*User, error)
	RegisterUser(email string) error
}

var (
	errUserInitiate = "user: initiate error"
)

type User struct {
	repo    user.UserRepository
	service *user.UserService

	ID    string
	Email string

	Config    *config.Config
	DebugLog  utility.LogFunc `inject:"true"`
	SystemLog utility.LogFunc `inject:"true"`
}
func NewUser(repo user.UserRepository, service *user.UserService) *User {
	return &User{
		repo:    repo,
		service: service,
	}
}
func (u *User) Initiate(ctx context.Context) (newCtx context.Context, err error) {
	var configs map[string]config.Option = map[string]config.Option{
		"QUEUE_PROVIDER":          config.Option{Type: config.String, Default: "kafka", EnvironmentKey: "QUEUE_PROVIDER"},
		"QUEUE_KAFKA_BROKERS":     config.Option{Type: config.String, EnvironmentKey: "QUEUE_KAFKA_BROKERS"},
		"QUEUE_KAFKA_EVENT_GROUP": config.Option{Type: config.String, EnvironmentKey: "QUEUE_KAFKA_EVENT_GROUP"},
		"QUEUE_KAFKA_EVENT_TOPIC": config.Option{Type: config.String, EnvironmentKey: "QUEUE_KAFKA_EVENT_TOPIC"},
	}
	newCtx, err = u.Config.SetOptions(ctx, configs)
	if err != nil {
		return newCtx, fmt.Errorf(errUserInitiate+": %s", err.Error())
	}
	return newCtx, nil
}
func (u *User) OnStartup(ctx context.Context) (context.Context, error) {
	if u.DebugLog == nil {
		u.DebugLog = utility.DefaultLogFunc
	} else {
		u.DebugLog = utility.DiscardLogFunc
	}
	if u.SystemLog == nil {
		u.SystemLog = utility.DefaultLogFunc
	} else {
		u.SystemLog = utility.DiscardLogFunc
	}
	return ctx, nil
}
func (u *User) OnShutdown(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
func (u *User) ListUser() ([]*User, error) {
	users, err := u.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return toUser(users), nil
}
func (u *User) RegisterUser(email string) error {
	uid := uuid.NewV4()
	if err := u.service.Duplicated(email); err != nil {
		return err
	}
	user := user.NewUser(uid.String(), email)
	if err := u.repo.Save(user); err != nil {
		return err
	}
	return nil
}
func toUser(users []*user.User) []*User {
	res := make([]*User, len(users))
	for i, user := range users {
		res[i] = &User{
			ID:    user.GetID(),
			Email: user.GetEmail(),
		}
	}
	return res
}
