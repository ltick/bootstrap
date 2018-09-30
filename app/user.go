package app

import (
	"context"
	"fmt"

	"github.com/ltick/dummy/app/domain/user"
	libConfig "github.com/ltick/tick-framework/module/config"
	libUtility "github.com/ltick/tick-framework/module/utility"
	"github.com/ltick/tick-routing"
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

	Config    *libConfig.Instance
	Utility   *libUtility.Instance
	DebugLog  libUtility.LogFunc `inject:"true"`
	SystemLog libUtility.LogFunc `inject:"true"`
}
func NewUser(repo user.UserRepository, service *user.UserService) *User {
	return &User{
		repo:    repo,
		service: service,
	}
}
func (this *User) Initiate(ctx context.Context) (newCtx context.Context, err error) {
	var configs map[string]libConfig.Option = map[string]libConfig.Option{
		"QUEUE_PROVIDER":          libConfig.Option{Type: libConfig.String, Default: "kafka", EnvironmentKey: "QUEUE_PROVIDER"},
		"QUEUE_KAFKA_BROKERS":     libConfig.Option{Type: libConfig.String, EnvironmentKey: "QUEUE_KAFKA_BROKERS"},
		"QUEUE_KAFKA_EVENT_GROUP": libConfig.Option{Type: libConfig.String, EnvironmentKey: "QUEUE_KAFKA_EVENT_GROUP"},
		"QUEUE_KAFKA_EVENT_TOPIC": libConfig.Option{Type: libConfig.String, EnvironmentKey: "QUEUE_KAFKA_EVENT_TOPIC"},
	}
	newCtx, err = this.Config.SetOptions(ctx, configs)
	if err != nil {
		return newCtx, fmt.Errorf(errUserInitiate+": %s", err.Error())
	}
	return newCtx, nil
}
func (this *User) OnStartup(ctx context.Context) (context.Context, error) {
	if this.DebugLog == nil {
		this.DebugLog = this.Utility.DefaultLogFunc
	}
	if this.SystemLog == nil {
		this.SystemLog = this.Utility.DefaultLogFunc
	}
	return ctx, nil
}
func (this *User) OnShutdown(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
func (this *User) OnRequestStartup(c *routing.Context) error {
	return nil
}
func (this *User) OnRequestShutdown(c *routing.Context) error {
	return nil
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
