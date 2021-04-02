package ginz

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/suisrc/config.zgo"
	"github.com/suisrc/logger.zgo"
)

// RunHTTPServer 初始化http服务
func RunHTTPServer(ctx context.Context, handler http.Handler) func() {
	conf := config.C.HTTP
	addr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)

	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	go func() {
		logger.Printf(ctx, "http server is running at %s.", addr)
		// var err error
		// if conf.CertFile != "" && conf.KeyFile != "" {
		// 	srv.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
		// 	err = srv.ListenAndServeTLS(conf.CertFile, conf.KeyFile)
		// } else {
		// 	err = srv.ListenAndServe()
		// }
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() {
		ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(conf.ShutdownTimeout))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			logger.Errorf(ctx, logger.ErrorWW(err))
		}
	}
}

// RunWithShutdown 运行服务
func RunWithShutdown(ctx context.Context, runServer func() (func(), error)) error {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	shutdownServer, err := runServer()
	if err != nil {
		return err
	}

	sig := <-sc // 等待服务器中断
	logger.Printf(ctx, "received a signal [%s]", sig.String())
	// 结束服务
	logger.Printf(ctx, "http server shutdown ...")
	shutdownServer()
	logger.Printf(ctx, "http server exiting")
	time.Sleep(time.Second) // 延迟1s, 用于日志等信息保存
	return nil
}
