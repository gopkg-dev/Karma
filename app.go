package karma

import (
	"context"
	"errors"
	"github.com/gopkg-dev/karma/log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
)

// App 是一个应用程序组件生命周期管理器。
type App struct {
	opts   options
	ctx    context.Context
	cancel func()
}

// New 创建一个应用程序生命周期管理器。
func New(opts ...Option) *App {
	o := options{
		ctx: context.Background(),
		// do not catch SIGKILL signal, need to waiting for kill self by others.
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
	for _, opt := range opts {
		opt(&o)
	}
	if o.logger != nil {
		log.SetLogger(o.logger)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{
		ctx:    ctx,
		cancel: cancel,
		opts:   o,
	}
}

// Name 返回服务名称。
func (a *App) Name() string { return a.opts.name }

// Version 返回应用版本。
func (a *App) Version() string { return a.opts.version }

// Run 按正确顺序执行所有钩子并启动所有服务器。
func (a *App) Run() error {
	eg, ctx := errgroup.WithContext(a.ctx)

	// 启动所有服务器
	wg := sync.WaitGroup{}
	for _, srv := range a.opts.servers {
		srv := srv
		eg.Go(func() error {
			<-ctx.Done() // 等待停止信号
			return srv.Stop(ctx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(ctx)
		})
	}

	// watch signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, a.opts.sigs...)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-quit:
				return a.Stop()
			}
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}

	return nil
}

// Stop 优雅地停止应用程序并执行 beforeStop 和 afterStop 钩子。
func (a *App) Stop() (err error) {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}
