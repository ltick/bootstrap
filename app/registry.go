package app

import (
	"github.com/ltick/dummy/app/domain/user"
	"github.com/ltick/dummy/app/interface/persistence/memory"
	"github.com/ltick/tick-framework"
)

var repo user.UserRepository = memory.NewUserRepository()
var service *user.UserService = user.NewUserService(repo)
var components []*ltick.Component = []*ltick.Component{
	// 存储初始化
	&ltick.Component{Name: "user", Component: NewUser(repo, service)},
}