package karma

import (
	"context"
	"os"

	"github.com/gopkg-dev/karma/log"
	"github.com/gopkg-dev/karma/transport"
)

// Option 是一个应用程序选项。
type Option func(o *options)

// options 是应用程序的选项。
type options struct {
	name    string
	version string

	ctx  context.Context
	sigs []os.Signal

	logger  log.Logger
	servers []transport.Server
}

// WithName 设置服务名称。
func WithName(name string) Option {
	return func(o *options) { o.name = name }
}

// WithVersion 设置服务版本。
func WithVersion(version string) Option {
	return func(o *options) { o.version = version }
}

// WithLogger 日志记录器
func WithLogger(logger log.Logger) Option {
	return func(o *options) { o.logger = logger }
}

// WithContext 设置服务上下文。
func WithContext(ctx context.Context) Option {
	return func(o *options) { o.ctx = ctx }
}

// WithServer 设置传输服务器。
func WithServer(srv ...transport.Server) Option {
	return func(o *options) { o.servers = srv }
}

// WithSignal 设置退出信号。
func WithSignal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}
