// Copyright © 2017 Ding Jing <dingjing@lianjia.com>
// this file is part of {{.appName}}

package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ltick/tick-framework"
	"github.com/ltick/tick-framework/module"
	libConfig "github.com/ltick/tick-framework/module/config"
	libLogger "github.com/ltick/tick-framework/module/logger"
	libUtility "github.com/ltick/tick-framework/module/utility"
	"github.com/ltick/tick-routing"
	"github.com/ltick/tick-routing/access"
	"github.com/spf13/cobra"

	appModule "github.com/ltick/dummy/app/module"
	appModule1 "github.com/ltick/dummy/app/module1"
)

var modules []*module.Module = []*module.Module{
	// 存储初始化
	&module.Module{Name: "appModule1", Module: &appModule1.Instance{}},
	&module.Module{Name: "appModule", Module: &appModule.Instance{}},
}
var configs map[string]libConfig.Option = map[string]libConfig.Option{
	"APP_ENV":     libConfig.Option{Type: libConfig.String, Default: "local", EnvironmentKey: "APP_ENV"},
	"PREFIX_PATH": libConfig.Option{Type: libConfig.String, Default: prefixPath, EnvironmentKey: "PREFIX_PATH"},
	"TMP_PATH":    libConfig.Option{Type: libConfig.String, Default: "/tmp"},
	"DEBUG":       libConfig.Option{Type: libConfig.String, Default: false},

	"ACCESS_LOG_TYPE":      libConfig.Option{Type: libConfig.String, Default: "console", EnvironmentKey: "ACCESS_LOG_TYPE"},
	"ACCESS_LOG_FILENAME":  libConfig.Option{Type: libConfig.String, Default: "/tmp/access.log", EnvironmentKey: "ACCESS_LOG_FILENAME"},
	"ACCESS_LOG_WRITER":    libConfig.Option{Type: libConfig.String, Default: "discard", EnvironmentKey: "ACCESS_LOG_WRITER"},
	"ACCESS_LOG_MAX_LEVEL": libConfig.Option{Type: libConfig.String, Default: libLogger.LevelInfo, EnvironmentKey: "ACCESS_LOG_MAX_LEVEL"},
	"ACCESS_LOG_FORMATTER": libConfig.Option{Type: libConfig.String, Default: "raw", EnvironmentKey: "ACCESS_LOG_FORMATTER"},

	"DEBUG_LOG_TYPE":      libConfig.Option{Type: libConfig.String, Default: "console", EnvironmentKey: "DEBUG_LOG_TYPE"},
	"DEBUG_LOG_FILENAME":  libConfig.Option{Type: libConfig.String, Default: "/tmp/debug.log", EnvironmentKey: "DEBUG_LOG_FILENAME"},
	"DEBUG_LOG_WRITER":    libConfig.Option{Type: libConfig.String, Default: "discard", EnvironmentKey: "DEBUG_LOG_WRITER"},
	"DEBUG_LOG_MAX_LEVEL": libConfig.Option{Type: libConfig.String, Default: libLogger.LevelInfo, EnvironmentKey: "DEBUG_LOG_MAX_LEVEL"},
	"DEBUG_LOG_FORMATTER": libConfig.Option{Type: libConfig.String, Default: "default", EnvironmentKey: "DEBUG_LOG_FORMATTER"},

	"SYSTEM_LOG_TYPE":      libConfig.Option{Type: libConfig.String, Default: "console", EnvironmentKey: "SYSTEM_LOG_TYPE"},
	"SYSTEM_LOG_FILENAME":  libConfig.Option{Type: libConfig.String, Default: "/tmp/system.log", EnvironmentKey: "SYSTEM_LOG_FILENAME"},
	"SYSTEM_LOG_WRITER":    libConfig.Option{Type: libConfig.String, Default: "discard", EnvironmentKey: "SYSTEM_LOG_WRITER"},
	"SYSTEM_LOG_MAX_LEVEL": libConfig.Option{Type: libConfig.String, Default: libLogger.LevelInfo, EnvironmentKey: "SYSTEM_LOG_MAX_LEVEL"},
	"SYSTEM_LOG_FORMATTER": libConfig.Option{Type: libConfig.String, Default: "sys", EnvironmentKey: "SYSTEM_LOG_FORMATTER"},
}

var enviroment string
var prefixPath string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "dummy-server [command]",
	Short: "Dummy Server",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "启动预处理",
	Long: `启动预处理

Example:
  dummy-server start
  dummy-server start --prefix $GOPATH/src/nebula/media
    `,
	Run: func(cmd *cobra.Command, args []string) {
		StartService()
	},
}

func DefaultLogFunc(ctx context.Context, format string, data ...interface{}) {
	forwardRequestId, requestId, _, serverAddress := GetLogContext(ctx)
	logData := make([]interface{}, len(data)+3)
	logData[0] = forwardRequestId
	logData[1] = requestId
	logData[2] = serverAddress
	copy(logData[3:], data)
	log.Printf("APP|%s|%s|%s|"+format, logData...)
}

func GetLogContext(ctx context.Context) (forwardRequestId string, requestId string, clientIP string, serverAddress string) {
	if ctx.Value("forwardRequestId") != nil {
		forwardRequestId = ctx.Value("forwardRequestId").(string)
	}
	if ctx.Value("requestId") != nil {
		requestId = ctx.Value("requestId").(string)
	}
	if ctx.Value("clientIP") != nil {
		clientIP = ctx.Value("clientIP").(string)
	}
	if ctx.Value("serverAddress") != nil {
		serverAddress = ctx.Value("serverAddress").(string)
	}
	return forwardRequestId, requestId, clientIP, serverAddress
}
func StartService() {
	if prefixPath == "" {
		fmt.Printf("dummy-server: prefix path does not set\n")
		os.Exit(1)
	}
	e := ltick.NewClassic(modules, configs, &ltick.Option{
		PathPrefix: prefixPath,
		EnvPrefix:  "DUMMY",
	})
	err := e.UseModule("cache", "queue", "database")
	if err != nil {
		e.SystemLog(fmt.Sprintf("dummy-server: use module error: " + err.Error()))
		os.Exit(1)
	}
	systemLogger, err := e.GetLogger("system")
	if err != nil {
		e.SystemLog(fmt.Sprintf("dummy-server: get logger error: " + err.Error()))
		os.Exit(1)
	}
	systemLogFunc := func(ctx context.Context, format string, data ...interface{}) {
		if systemLogger == nil {
			return
		}
		systemLogger.Info(format, data...)
	}
	debugLogger, err := e.GetLogger("debug")
	if err != nil {
		e.SystemLog(fmt.Sprintf("dummy-server: get logger error: " + err.Error()))
		os.Exit(1)
	}
	debugLogFunc := func(ctx context.Context, format string, data ...interface{}) {
		if debugLogger == nil {
			return
		}
		forwardRequestId, requestId, _, serverAddress := GetLogContext(ctx)
		logData := make([]interface{}, len(data)+3)
		logData[0] = forwardRequestId
		logData[1] = requestId
		logData[2] = serverAddress
		copy(logData[3:], data)
		debugLogger.Debug("APP|%s|%s|%s|"+format, logData...)
	}
	accessLogger, err := e.GetLogger("access")
	if err != nil {
		e.SystemLog(fmt.Sprintf("dummy-server: get logger error: " + err.Error()))
		os.Exit(1)
	}
	traceLogFunc := func(ctx context.Context, format string, data ...interface{}) {
		if debugLogger == nil {
			return
		}
		forwardRequestId, requestId, _, serverAddress := GetLogContext(ctx)
		logData := make([]interface{}, len(data)+3)
		logData[0] = forwardRequestId
		logData[1] = requestId
		logData[2] = serverAddress
		copy(logData[3:], data)
		debugLogger.Info("TRACE|%s|%s|%s|"+format, logData...)
	}
	accessLogFunc := func(c *routing.Context, rw *access.LogResponseWriter, elapsed float64) {
		if debugLogger == nil || accessLogger == nil {
			return
		}
		forwardRequestId, requestId, clientIP, serverAddress := GetLogContext(c.Context)
		requestLine := fmt.Sprintf("%s %s %s", c.Request.Method, c.Request.RequestURI, c.Request.Proto)
		debug := new(bool)
		if c.Context.Value("DEBUG") != nil {
			*debug = c.Context.Value("DEBUG").(bool)
		}
		if *debug {
			debugLogger.Info(`ACCESS|%s|%s|%s - %s [%s] "%s" %d %d %d %.3f "%s" "%s" %s %s "%v" "%v"`, forwardRequestId, requestId, clientIP, c.Request.Host, time.Now().Format("2/Jan/2006:15:04:05 -0700"), requestLine, c.Request.ContentLength, rw.Status, rw.BytesWritten, elapsed/1e3, c.Request.Header.Get("Referer"), c.Request.Header.Get("User-Agent"), c.Request.RemoteAddr, serverAddress, c.Request.Header, rw.Header())
		} else {
			debugLogger.Info(`ACCESS|%s|%s|%s - %s [%s] "%s" %d %d %d %.3f "%s" "%s" %s %s "-" "-"`, forwardRequestId, requestId, clientIP, c.Request.Host, time.Now().Format("2/Jan/2006:15:04:05 -0700"), requestLine, c.Request.ContentLength, rw.Status, rw.BytesWritten, elapsed/1e3, c.Request.Header.Get("Referer"), c.Request.Header.Get("User-Agent"), c.Request.RemoteAddr, serverAddress)
		}
		if *debug {
			accessLogger.Info(`%s - %s [%s] "%s" %d %d %d %.3f "%s" "%s" %s %s "%v" "%v"`, clientIP, c.Request.Host, time.Now().Format("2/Jan/2006:15:04:05 -0700"), requestLine, c.Request.ContentLength, rw.Status, rw.BytesWritten, elapsed/1e3, c.Request.Header.Get("Referer"), c.Request.Header.Get("User-Agent"), c.Request.RemoteAddr, serverAddress, c.Request.Header, rw.Header())
		} else {
			accessLogger.Info(`%s - %s [%s] "%s" %d %d %d %.3f "%s" "%s" %s %s "-" "-"`, clientIP, c.Request.Host, time.Now().Format("2/Jan/2006:15:04:05 -0700"), requestLine, c.Request.ContentLength, rw.Status, rw.BytesWritten, elapsed/1e3, c.Request.Header.Get("Referer"), c.Request.Header.Get("User-Agent"), c.Request.RemoteAddr, serverAddress)
		}
	}
	recoveryHandler := func(c *routing.Context, err error) error {
		if debugLogger == nil {
			return routing.NewHTTPError(http.StatusInternalServerError)
		}
		forwardRequestId, requestId, _, serverAddress := GetLogContext(c.Context)
		if httpError, ok := err.(routing.HTTPError); ok {
			status := httpError.StatusCode()
			switch status {
			case http.StatusBadRequest:
				fallthrough
			case http.StatusForbidden:
				fallthrough
			case http.StatusNotFound:
				fallthrough
			case http.StatusRequestTimeout:
				fallthrough
			case http.StatusMethodNotAllowed:
				debugLogger.Info("CLIENT_ERROR|%s|%s|%s|%s", forwardRequestId, requestId, serverAddress, httpError.Error())
			default:
				debugLogger.Error("SERVER_ERROR|%s|%s|%s|%s", forwardRequestId, requestId, serverAddress, httpError.Error())
			}
			return routing.NewHTTPError(status, httpError.Error())
		} else {
			debugLogger.Emergency("SERVER_ERROR|%s|%s|%s|%s", forwardRequestId, requestId, serverAddress, err.Error())
			return routing.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}
	var values map[string]interface{} = map[string]interface{}{
		"DebugLog":  debugLogFunc,
		"TraceLog":  traceLogFunc,
		"SystemLog": systemLogFunc,
	}
	e.WithValues(values).NewClassicServer("test", func(c *routing.Context) error {
		return nil
	}).SetLogFunc(accessLogFunc, systemLogger.Emergency, recoveryHandler)
	err = e.Startup()
	if err != nil {
		e.SystemLog(fmt.Sprintf("dummy-server: startup error: " + err.Error()))
		os.Exit(1)
	}
	e.ListenAndServe()
	err = e.Shutdown()
	if err != nil {
		e.SystemLog(fmt.Sprintf("dummy-server: shutdown error: " + err.Error()))
		os.Exit(1)
	}
}

type Callback struct {
	Utility *libUtility.Instance
	Config  *libConfig.Instance
}

func (f *Callback) OnStartup(a *ltick.Engine) error {
	a.SystemLog("dummy-server: application startup")
	return nil
}

func (f *Callback) OnShutdown(a *ltick.Engine) error {
	a.SystemLog("dummy-server: application shutdown")
	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	executeFile, _ := exec.LookPath(os.Args[0])
	executePath, _ := filepath.Abs(executeFile)
	defaultPrefixPath := filepath.Dir(filepath.Dir(executePath))
	RootCmd.PersistentFlags().StringVar(&prefixPath, "prefix", defaultPrefixPath, "prefix path of this service")
	RootCmd.AddCommand(startCmd)
	if err := RootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
