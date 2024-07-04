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

// BodyLimit with app Max body size that the server accepts.
func BodyLimit(v int) Option {
	return func(o *Server) { o.bodyLimit = v }
}

// AppName with server app name.
func AppName(v string) Option {
	return func(o *Server) { o.appName = v }
}

func ServerHeader(v string) Option {
	return func(o *Server) { o.serverHeader = v }
}

func Concurrency(v int) Option {
	return func(o *Server) { o.concurrency = v }
}

func DisableKeepalive(v bool) Option {
	return func(o *Server) { o.disableKeepalive = v }
}

func EnablePrintRoutes(v bool) Option {
	return func(o *Server) { o.enablePrintRoutes = v }
}

func IdleTimeout(second int) Option {
	return func(o *Server) { o.idleTimeout = second }
}

func ReadTimeout(second int) Option {
	return func(o *Server) { o.readTimeout = second }
}

func WriteTimeout(second int) Option {
	return func(o *Server) { o.writeTimeout = second }
}

func ShutdownTimeout(second int) Option {
	return func(o *Server) { o.shutdownTimeout = second }
}

func ErrorHandler(t fiber.ErrorHandler) Option {
	return func(o *Server) { o.defaultErrorHandler = t }
}

func JSONEncoder(v utils.JSONMarshal) Option {
	return func(o *Server) { o.jsonEncoder = v }
}

func JSONDecoder(v utils.JSONUnmarshal) Option {
	return func(o *Server) { o.jsonDecoder = v }
}

// Middleware with server middleware option.
func Middleware(m ...fiber.Handler) Option {
	return func(o *Server) {
		o.handlers = m
	}
}
