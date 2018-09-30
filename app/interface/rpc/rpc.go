package rpc

import (
	"github.com/ltick/dummy/app/interface/rpc/v1.0"
	"github.com/ltick/tick-framework"
	"google.golang.org/grpc"
)

func Apply(server *grpc.Server, engine *ltick.Engine) {
	v1.Apply(server, engine)
}