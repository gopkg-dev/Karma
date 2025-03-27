package fiberx

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// Option is config option.
type Option func(app *Server)

func ServerHost(host string) Option {
	return func(s *Server) {
		s.host = host
	}
}

func ServerPort(port int) Option {
	return func(s *Server) {
		s.port = port
	}
}

// BodyLimit 设置请求主体允许的最大大小，如果大小超出配置的限制，则会发送 413 - 请求实体太大响应。
func BodyLimit(v int) Option {
	return func(o *Server) { o.bodyLimit = v }
}

// AppName 应用程序名称
func AppName(v string) Option {
	return func(o *Server) { o.appName = v }
}

// ServerHeader 设置服务器标头
func ServerHeader(v string) Option {
	return func(o *Server) { o.serverHeader = v }
}

// Concurrency 设置服务器的并发连接数
func Concurrency(v int) Option {
	return func(o *Server) { o.concurrency = v }
}

// CaseSensitive 设置路由是否区分大小写
func CaseSensitive(v bool) Option {
	return func(o *Server) { o.caseSensitive = v }
}

// StrictRouting 设置路由是否严格匹配
func StrictRouting(v bool) Option {
	return func(o *Server) { o.strictRouting = v }
}

// StreamRequestBody 设置是否启用流式请求主体
func StreamRequestBody(v bool) Option {
	return func(o *Server) { o.streamRequestBody = v }
}

// DisablePreParseMultipartForm 禁用在路由处理程序之前解析多部分表单
func DisablePreParseMultipartForm(v bool) Option {
	return func(o *Server) { o.disablePreParseMultipartForm = v }
}

// DisableKeepalive 禁用保持连接，服务器将在向客户端发送第一个响应后关闭传入连接
func DisableKeepalive(v bool) Option {
	return func(o *Server) { o.disableKeepalive = v }
}

// EnablePrintRoutes 启用打印路由
func EnablePrintRoutes(v bool) Option {
	return func(o *Server) { o.enablePrintRoutes = v }
}

// IdleTimeout 设置服务器等待所有连接关闭的时间
func IdleTimeout(second int) Option {
	return func(o *Server) { o.idleTimeout = second }
}

// ReadTimeout 设置服务器读取请求头和主体的时间
func ReadTimeout(second int) Option {
	return func(o *Server) { o.readTimeout = second }
}

// WriteTimeout 设置服务器写入响应的时间
func WriteTimeout(second int) Option {
	return func(o *Server) { o.writeTimeout = second }
}

func WriteBufferSize(v int) Option {
	return func(o *Server) { o.writeBufferSize = v }

}

// ShutdownTimeout 设置服务器关闭的时间
func ShutdownTimeout(second int) Option {
	return func(o *Server) { o.shutdownTimeout = second }
}

// ErrorHandler 设置错误处理程序
func ErrorHandler(t fiber.ErrorHandler) Option {
	return func(o *Server) { o.defaultErrorHandler = t }
}

// JSONEncoder 设置 JSON 编码器
func JSONEncoder(v utils.JSONMarshal) Option {
	return func(o *Server) { o.jsonEncoder = v }
}

// JSONDecoder 设置 JSON 解码器
func JSONDecoder(v utils.JSONUnmarshal) Option {
	return func(o *Server) { o.jsonDecoder = v }
}

// Middleware with server middleware option.
func Middleware(m ...fiber.Handler) Option {
	return func(o *Server) {
		o.handlers = m
	}
}
