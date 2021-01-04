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

// DefAutoCompleteHandler name of DI definition
const DefAutoCompleteHandler = "http.handler.autocomplete"

// Definition init func.
func init() {
	container.Register(func(builder *container.ProviderObject) error {
		return builder.Add(container.Definition{
			Name: DefAutoCompleteHandler,
			Build: func(
				log logger.Logger, searchRepo domain.ISearchRepository, calendarRepo domain.ICalendarRepository,
			) (_ echo.Handler, err error) {
				return handlers.NewSearchHandler(
					log,
					validator.New(),
					searchRepo,
					calendarRepo,
				), nil
			},
			Params: container.Params{
				"0": container.Service(logger.DefLogger),
				"1": container.Service(postgres.DefSearchPostgresRepo),
				"2": container.Service(postgres.DefCalendarPostgresRepo),
			},
		})
	})
}
