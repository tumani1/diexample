// Package echo provide dependency injection definitions.
package echo

import (
	"net/http"

	"github.com/labstack/echo"

	"github.com/tumani1/diexample/di/sarulabsdi/container"
	"github.com/tumani1/diexample/di/sarulabsdi/echo/errors"
)

// DefErrorHandler definition name.
const DefErrorHandler = "echo.error_handler"

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder) error {
		return builder.Add(container.Definition{
			Name: DefErrorHandler,
			Build: func(ctx container.Context) (_ interface{}, err error) {
				return func(err error, c echo.Context) {
					var (
						e      = c.Echo()
						logger = c.Logger()
						code   = http.StatusInternalServerError
						msg    = http.StatusText(code)
						errMsg interface{}
					)

					switch he := err.(type) {
					case *echo.HTTPError:
						code = he.Code
						msg = he.Message.(string)
					case *errors.HTTPError:
						code = he.Code
						msg = he.Message
						errMsg = he.Description
					}

					if e.Debug {
						msg = err.Error()
					}

					logger.Error(err)

					if !c.Response().Committed {
						if c.Request().Method == echo.HEAD {
							err = c.NoContent(code)
						} else {
							var m = echo.Map{"message": msg}
							if errMsg != nil {
								m["error"] = errMsg
							}
							err = c.JSON(code, m)
						}
						if err != nil {
							logger.Error(err)
						}
					}
				}, nil
			},
		})
	})
}
