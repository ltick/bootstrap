package v1

import (
	"github.com/ltick/dummy/app/interface/rpc/v1.0/protocol"
	"github.com/ltick/tick-framework"
	"github.com/ltick/dummy/app"
	"google.golang.org/grpc"
)

func Apply(server *grpc.Server, engine *ltick.Engine) error {
	user, err := engine.GetModule("user")
	if err != nil {
		return err
	}
	protocol.RegisterUserServiceServer(server, NewUserService(user.(app.User)))
	return nil
}