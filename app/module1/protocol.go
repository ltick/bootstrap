package module

import (
	"context"
	"fmt"

	appModule "github.com/ltick/dummy/app/module"
	libConfig "github.com/ltick/tick-framework/module/config"
	libUtility "github.com/ltick/tick-framework/module/utility"
	"github.com/ltick/tick-routing"
)

var (
	errInitiate = "test module: initiate error"
	errStartup  = "test module: startup error"
)

var debugLog libUtility.LogFunc
var systemLog libUtility.LogFunc

type Instance struct {
	AppModule *appModule.Instance `inject:"true"`
	Config    *libConfig.Instance
	Utility   *libUtility.Instance
	DebugLog  libUtility.LogFunc `inject:"true"`
	SystemLog libUtility.LogFunc `inject:"true"`
}

func NewInstance() *Instance {
	return &Instance{}
}
func (this *Instance) Initiate(ctx context.Context) (newCtx context.Context, err error) {
	var configs map[string]libConfig.Option = map[string]libConfig.Option{
		"QUEUE_PROVIDER":          libConfig.Option{Type: libConfig.String, Default: "kafka", EnvironmentKey: "QUEUE_PROVIDER"},
		"QUEUE_KAFKA_BROKERS":     libConfig.Option{Type: libConfig.String, EnvironmentKey: "QUEUE_KAFKA_BROKERS"},
		"QUEUE_KAFKA_EVENT_GROUP": libConfig.Option{Type: libConfig.String, EnvironmentKey: "QUEUE_KAFKA_EVENT_GROUP"},
		"QUEUE_KAFKA_EVENT_TOPIC": libConfig.Option{Type: libConfig.String, EnvironmentKey: "QUEUE_KAFKA_EVENT_TOPIC"},
	}
	newCtx, err = this.Config.SetOptions(ctx, configs)
	if err != nil {
		return newCtx, fmt.Errorf(errInitiate+": %s", err.Error())
	}
	return newCtx, nil
}
func (this *Instance) OnStartup(ctx context.Context) (context.Context, error) {
	if this.DebugLog != nil {
		debugLog = this.DebugLog
	} else {
		debugLog = this.Utility.DefaultLogFunc
	}
	if this.SystemLog != nil {
		systemLog = this.SystemLog
	} else {
		systemLog = this.Utility.DefaultLogFunc
	}
	return ctx, nil
}
func (this *Instance) OnShutdown(ctx context.Context) (context.Context, error) {
	return ctx, nil
}
func (this *Instance) OnRequestStartup(c *routing.Context) error {
	return nil
}
func (this *Instance) OnRequestShutdown(c *routing.Context) error {
	return nil
}
