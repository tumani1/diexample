package handlers

import (
	"gopkg.in/go-playground/validator.v9"

	"github.com/tumani1/diexample/di/sarulabsdingo/container"
	"github.com/tumani1/diexample/di/sarulabsdingo/definition/echo"
	"github.com/tumani1/diexample/di/sarulabsdingo/definition/logger"
	"github.com/tumani1/diexample/di/sarulabsdingo/internal/definition/postgres"
	"github.com/tumani1/diexample/di/sarulabsdingo/internal/domain"
	"github.com/tumani1/diexample/di/sarulabsdingo/internal/http/handlers"
)

// DefSearchHandler name of DI definition
const DefSearchHandler = "http.handler.search"

// Definition init func.
func init() {
	container.Register(func(builder *container.ProviderObject) error {
		return builder.Add(container.Definition{
			Name: DefSearchHandler,
			Build: func(log logger.Logger, autoCompleteRepo domain.IAutoCompleteRepository) (_ echo.Handler, err error) {
				return handlers.NewAutoCompleteHandler(
					log,
					validator.New(),
					autoCompleteRepo,
				), nil
			},
			Params: container.Params{
				"0": container.Service(logger.DefLogger),
				"1": container.Service(postgres.DefAutoCompletePostgresRepo),
			},
		})
	})
}
