// Package echo provide dependency injection definitions.
package echo

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"

	"github.com/tumani1/diexample/di/sarulabsdingo/container"
	"github.com/tumani1/diexample/di/sarulabsdingo/definition/config"
	appEcho "github.com/tumani1/diexample/di/sarulabsdingo/echo"
)

const (
	// DefEcho definition name.
	DefEcho = "echo"

	// DefHTTPHandlerTag http handler tag name.
	DefHTTPHandlerTag = "http.handler"
)

type (
	// Echo type alias of *echo.Echo
	Echo = *echo.Echo

	// Handler type alias of http.Handler
	Handler = appEcho.Handler
)

func init() {
	container.Register(func(builder *container.ProviderObject) error {
		return builder.Add(container.Definition{
			Name: DefEcho,
			Build: func(cfg config.Config, errHandler func(err error, c echo.Context)) (_ Echo, err error) {
				var e = echo.New()
				e.Debug = cfg.GetBool("http.debug")
				e.HTTPErrorHandler = errHandler

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

				e.Use(middleware.Recover())

				return e, nil
			},
			Params: container.Params{
				"0": container.Service(config.DefConfig),
			},
		})
	})
}
