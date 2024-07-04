package fiberx

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"sort"
	"text/tabwriter"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
)

type Server struct {
	*fiber.App
	host                   string
	port                   int
	appName                string
	serverHeader           string
	concurrency            int
	bodyLimit              int
	caseSensitive          bool
	strictRouting          bool
	disableStartupMessage  bool
	disableKeepalive       bool
	enablePrintRoutes      bool
	views                  fiber.Views
	viewsLayout            string
	idleTimeout            int
	readTimeout            int
	writeTimeout           int
	shutdownTimeout        int
	jsonEncoder            utils.JSONMarshal
	jsonDecoder            utils.JSONUnmarshal
	handlers               []fiber.Handler
	defaultErrorHandler    fiber.ErrorHandler
	defaultNotFoundHandler fiber.Handler
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		appName:                "go-fiber-admin",
		serverHeader:           "go-fiber-server",
		concurrency:            256 * 1024,
		bodyLimit:              5, // MB
		caseSensitive:          true,
		strictRouting:          true,
		disableStartupMessage:  true,
		disableKeepalive:       false,
		enablePrintRoutes:      false,
		idleTimeout:            10,
		readTimeout:            60,
		writeTimeout:           60,
		shutdownTimeout:        10,
		jsonEncoder:            json.Marshal,
		jsonDecoder:            json.Unmarshal,
		defaultErrorHandler:    DefaultErrorHandler,
		defaultNotFoundHandler: DefaultNotFoundHandler,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.App = fiber.New(fiber.Config{
		Prefork:               false,
		ServerHeader:          s.serverHeader,
		StrictRouting:         s.strictRouting,
		CaseSensitive:         s.caseSensitive,
		BodyLimit:             s.bodyLimit * 1024 * 1024,
		Concurrency:           s.concurrency,
		ReadTimeout:           time.Duration(s.readTimeout) * time.Second,
		WriteTimeout:          time.Duration(s.writeTimeout) * time.Second,
		IdleTimeout:           time.Duration(s.idleTimeout) * time.Second,
		ErrorHandler:          s.defaultErrorHandler,
		DisableKeepalive:      s.disableKeepalive,
		DisableStartupMessage: s.disableStartupMessage,
		//EnablePrintRoutes:     s.enablePrintRoutes,
		AppName:     s.appName,
		JSONEncoder: s.jsonEncoder,
		JSONDecoder: s.jsonDecoder,
		Views:       s.views,
		ViewsLayout: s.viewsLayout,
	})

	// register middleware
	if s.handlers != nil {
		for _, handler := range s.handlers {
			s.App.Use(handler)
		}
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	// Prepare an endpoint for 'Not Found'.
	s.Use(s.defaultNotFoundHandler)

	if s.enablePrintRoutes {
		s.PrintRoutes()
	}

	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	fmt.Printf("http server listening on %s\n", addr)

	err := s.Listen(addr)
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(s.shutdownTimeout))
	defer cancel()
	fmt.Printf("http server shutting down\n")
	return s.App.ShutdownWithContext(ctx)
}

func (s *Server) PrintRoutes() {
	colors := s.Config().ColorScheme

	out := colorable.NewColorableStdout()
	if os.Getenv("TERM") == "dumb" || os.Getenv("NO_COLOR") == "1" || (!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd())) {
		out = colorable.NewNonColorable(os.Stdout)
	}

	w := tabwriter.NewWriter(out, 1, 1, 1, ' ', 0)

	_, _ = fmt.Fprintf(w, "%sMethod\t%s| %sPath\t%s| %sName\t%s| %sHandlers\t%s\n", colors.Blue, colors.White, colors.Green, colors.White, colors.Cyan, colors.White, colors.Yellow, colors.Reset)
	_, _ = fmt.Fprintf(w, "%s------\t%s| %s----\t%s| %s----\t%s| %s--------\t%s\n", colors.Blue, colors.White, colors.Green, colors.White, colors.Cyan, colors.White, colors.Yellow, colors.Reset)
	for _, route := range s.GetRoutes(true) {
		_, _ = fmt.Fprintf(w, "%s%s\t%s| %s%s\t%s| %s%s\t%s| %s%s%s(%d handlers)\n", colors.Blue, route.Method, colors.White, colors.Green, route.Path, colors.White, colors.Cyan, route.Name, colors.White, colors.Yellow, route.HandlerName, colors.Reset, route.HandlersCount)
	}

	_ = w.Flush() //nolint:errcheck // It is fine to ignore the error here
}

// Route is a struct that holds all metadata for each registered handler.
type Route struct {
	Method        string   `json:"method"` // HTTP method
	Name          string   `json:"name"`   // Route's name
	Path          string   `json:"path"`   // Original registered route path
	Params        []string `json:"params"` // Case sensitive param keys
	HandlerName   string   `json:"-"`      //
	HandlersCount int      `json:"-"`      //
}

func (s *Server) GetRoutes(filter bool) []Route {
	var rs []Route
	for _, route := range s.App.GetRoutes(true) {
		if filter && route.Method == fiber.MethodHead {
			continue
		}
		nuHandlers := len(route.Handlers)
		lastHandler := route.Handlers[len(route.Handlers)-1]
		handlerName := runtime.FuncForPC(reflect.ValueOf(lastHandler).Pointer()).Name()
		routeName := route.Name
		if routeName == "" {
			routeName = "<No Name>"
		}
		rs = append(rs, Route{
			Method:        route.Method,
			Name:          routeName,
			Path:          route.Path,
			Params:        route.Params,
			HandlerName:   handlerName,
			HandlersCount: nuHandlers,
		})
	}
	sort.Slice(rs, func(i, j int) bool {
		return rs[i].Path < rs[j].Path
	})
	return rs
}
