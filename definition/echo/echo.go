// Package echo provide dependency injection definitions.
package echo

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"

	"gitlab.com/igor.tumanov1/theboatscom/container"
	"gitlab.com/igor.tumanov1/theboatscom/definition/config"
	appEcho "gitlab.com/igor.tumanov1/theboatscom/echo"
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
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
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

				if err = container.Iterate(ctx, DefHTTPHandlerTag, func(ctx container.Context, tag *container.Tag, name string) (err error) {
					var c appEcho.Handler
					if err = ctx.Fill(name, &c); err != nil {
						return err
					}

					c.Serve(e)
					return nil
				}); err != nil {
					return nil, errors.Wrap(err, "can't serve handler from container")
				}

				e.Use(middleware.Recover())

				return e, nil
			},
		})
	})
}
