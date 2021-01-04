// Package echo provide dependency injection definitions.
package echo

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"

	"github.com/tumani1/diexample/di/sarulabsdi/container"
	"github.com/tumani1/diexample/di/sarulabsdi/definition/config"
	appEcho "github.com/tumani1/diexample/di/sarulabsdi/echo"
)

const (
	// DefEcho definition name.
	DefEcho = "echo"

	// DefHTTPHandlerTag http handler tag name.
	DefHTTPHandlerTag = "http.handler.tag"
)

type (
	// Echo type alias of *echo.Echo
	Echo = *echo.Echo

	// Handler type alias of http.Handler
	Handler = appEcho.Handler
)

func init() {
	container.Register(func(builder *container.Builder) error {
		return builder.Add(container.Definition{
			Name: DefEcho,
			Build: func(ctx container.Context) (_ interface{}, err error) {
				var cfg config.Config
				if err = ctx.Fill(config.DefConfig, &cfg); err != nil {
					return nil, errors.Wrap(err, "can't get config from container")
				}

				var e = echo.New()
				e.Debug = cfg.GetBool("http.debug")

				switch cfg.GetString("http.level") {
				case "debug":
					e.Logger.SetLevel(log.DEBUG)
				case "info":
					e.Logger.SetLevel(log.INFO)
				case "warn":
					e.Logger.SetLevel(log.WARN)
				case "error":
					e.Logger.SetLevel(log.ERROR)
				case "off":
					e.Logger.SetLevel(log.OFF)
				}

				if err = ctx.Fill(DefErrorHandler, &e.HTTPErrorHandler); err != nil {
					return nil, errors.Wrap(err, "can't get error handler from container")
				}

				e.Use(middleware.Recover())

				return e, nil
			},
		})
	})
}
