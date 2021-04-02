package ginz

import (
	"context"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/config.zgo"
	"github.com/suisrc/logger.zgo"
)

//====================================
// 默认启动程序, 可以直接重新替换
//====================================

// Options options
type Options struct {
	ConfigFile    string
	Version       string
	BuildInjector BuildInjector
}

// Option 定义配置项
type Option func(*Options)

// BuildInjector 构建注入器的方法
type BuildInjector func() (engine *gin.Engine, clean func(), err error)

// Run 运行服务, 注意,必须对BuildInjector进行初始化
func Run(ctx context.Context, opts ...Option) error {
	return RunWithShutdown(ctx, func() (func(), error) {
		return RunServer(ctx, opts...)
	})
}

// RunServer 启动服务
func RunServer(ctx context.Context, opts ...Option) (func(), error) {
	var o Options
	for _, opt := range opts {
		opt(&o)
	}
	logger.SetVersion(o.Version)
	// 加载配置文件
	config.MustLoad(o.ConfigFile)
	config.Print()

	// 启动日志
	logger.Printf(ctx, "http server startup, M[%s]-V[%s]-P[%d]", config.C.RunMode, o.Version, os.Getpid())

	// 初始化日志模块
	loggerCleanFunc, err := logger.InitLogger(ctx)
	if err != nil {
		return nil, err
	}

	// 初始化依赖注入器
	engine, injectorCleanFunc, err := o.BuildInjector()
	if err != nil {
		return nil, err
	}

	// 初始化HTTP服务
	shutdownServerFunc := RunHTTPServer(ctx, engine)

	return func() {
		shutdownServerFunc()
		injectorCleanFunc()
		loggerCleanFunc()
	}, nil
}
