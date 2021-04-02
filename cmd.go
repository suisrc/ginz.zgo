package ginz

import (
	"context"
	"os"

	"github.com/suisrc/logger.zgo"
	"github.com/urfave/cli/v2"
)

// SetConfigFile 设定配置文件
func SetServeConfig(s string) Option {
	return func(o *Options) {
		o.ConfigFile = s
	}
}

// SetVersion 设定版本号b
func SetServeVersion(s string) Option {
	return func(o *Options) {
		o.Version = s
	}
}

// SetBuildInjector 设定注入助手
func SetBuildInjector(b BuildInjector) Option {
	return func(o *Options) {
		o.BuildInjector = b
	}
}

// CmdWeb ...
func CmdWeb(ctx context.Context, action func(c *cli.Context) error) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "运行web服务",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "conf",
				Aliases:     []string{"c"},
				Usage:       "配置文件(.json,.yaml,.toml)",
				DefaultText: "config.toml",
				//Required:   true,
			},
		},
		Action: action,
	}
}

func RunWeb(app *cli.App, run func(ctx context.Context, opts ...Option) error) {
	ctx := context.Background()
	app.Commands = []*cli.Command{
		CmdWeb(ctx, func(c *cli.Context) error {
			sc := SetServeConfig(c.String("conf"))
			sv := SetServeVersion(app.Version)
			return run(ctx, sc, sv)
		}),
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Errorf(ctx, logger.ErrorWW(err))
	}
}
