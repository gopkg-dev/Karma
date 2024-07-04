package logger

import (
	"fmt"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gopkg-dev/karma/errors"
)

func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Set PID once
	pid := os.Getpid()

	// Set variables
	var (
		once       sync.Once
		errHandler fiber.ErrorHandler
	)

	// Return new handler
	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Set error handler once
		once.Do(func() {
			// override error handler
			errHandler = c.App().ErrorHandler
		})

		// Set latency start time
		start := time.Now()

		// Handle request, store err for logging
		chainErr := c.Next()

		// Manually call error handler
		if chainErr != nil {
			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError) // TODO: Explain why we ignore the error here
			}
		}

		// Set latency stop time
		latency := time.Since(start)

		// Create the attrs slice
		attrs := make([]slog.Attr, 0)

		method := c.Method()
		reqBody := c.Body()
		resBody := c.Response().Body()
		statusCode := c.Response().StatusCode()
		contentType := c.Get(fiber.HeaderContentType)

		for _, key := range cfg.BuiltinAttrs {
			switch key {
			case TagReferer:
				attrs = append(attrs, slog.String(key, c.Get(fiber.HeaderReferer)))
			case TagProtocol:
				attrs = append(attrs, slog.String(key, c.Protocol()))
			case TagPort:
				attrs = append(attrs, slog.String(key, c.Port()))
			case TagIP:
				attrs = append(attrs, slog.String(key, c.IP()))
			case TagIPs:
				attrs = append(attrs, slog.String(key, c.Get(fiber.HeaderXForwardedFor)))
			case TagHost:
				attrs = append(attrs, slog.String(key, c.Hostname()))
			case TagPath:
				attrs = append(attrs, slog.String(key, c.Path()))
			case TagUrl:
				attrs = append(attrs, slog.String(key, c.OriginalURL()))
			case TagUA:
				attrs = append(attrs, slog.String(key, c.Get(fiber.HeaderUserAgent)))
			case TagBody:
				attr := slog.String(key, "N/A") // Set default value
				if method == fiber.MethodPost || method == fiber.MethodPut {
					mediaType, _, _ := mime.ParseMediaType(contentType)
					if mediaType == "application/json" {
						if len(reqBody) <= cfg.MaxRequestBodyLen {
							attrs = append(attrs, slog.String(key, string(reqBody)))
						} else {
							attrs = append(attrs, slog.String(key, "Exceeded size limit"))
						}
					}
				}
				attrs = append(attrs, attr)
			case TagBytesReceived:
				attrs = append(attrs, slog.Int(key, len(reqBody)))
			case TagBytesSent:
				attrs = append(attrs, slog.Int(key, len(resBody)))
			case TagRoute:
				attrs = append(attrs, slog.String(key, c.Route().Path))
			case TagResBody:
				if len(resBody) <= cfg.MaxResponseBodyLen {
					attrs = append(attrs, slog.String(key, string(resBody)))
				} else {
					attrs = append(attrs, slog.String(key, "Exceeded size limit"))
				}
			case TagReqHeaders:
				headersAttrs := make([]slog.Attr, 0)
				for k, v := range c.GetReqHeaders() {
					if _, found := HiddenRequestHeaders[strings.ToLower(k)]; found {
						headersAttrs = append(headersAttrs, slog.String(k, "*"))
						continue
					}
					headersAttrs = append(headersAttrs, slog.String(k, strings.Join(v, ",")))
				}
				attrs = append(attrs, slog.Any(TagReqHeaders, headersAttrs))
			case TagQueryStringParams:
				attrs = append(attrs, slog.String(key, c.Request().URI().QueryArgs().String()))
			case TagStatus:
				attrs = append(attrs, slog.Int(key, statusCode))
			case TagMethod:
				attrs = append(attrs, slog.String(key, c.Method()))
			case TagPid:
				attrs = append(attrs, slog.Int(key, pid))
			case TagLatency:
				attrs = append(attrs, slog.Duration(key, latency.Round(time.Microsecond)))
			case TagContentType:
				attrs = append(attrs, slog.String(key, contentType))
			case TagError:
				if chainErr != nil {
					e := errors.FromError(chainErr)
					attrs = append(attrs, slog.Any(TagError, e))
				}
			default:
				panic(fmt.Sprintf("Unexpected key value: %s", key))
			}
		}

		if cfg.CustomAttr != nil {
			for _, fn := range cfg.CustomAttr {
				attrs = append(attrs, fn(c, chainErr)...)
			}
		}

		var logLevel slog.Level

		switch {
		case statusCode >= fiber.StatusBadRequest && statusCode < fiber.StatusInternalServerError:
			logLevel = slog.LevelWarn
		case statusCode >= http.StatusInternalServerError:
			logLevel = slog.LevelError
		default:
			logLevel = slog.LevelInfo
		}

		cfg.Logger.LogAttrs(c.UserContext(), logLevel, "request", attrs...)

		if cfg.Done != nil {
			attributes := make(fiber.Map, len(attrs))
			for _, item := range attrs {
				attributes[item.Key] = item.Value.Any()
			}
			cfg.Done(c, attributes)
		}

		return nil
	}
}
