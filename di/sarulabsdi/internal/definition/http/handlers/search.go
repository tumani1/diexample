package handlers

import (
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"

	"github.com/tumani1/diexample/di/sarulabsdi/container"
	"github.com/tumani1/diexample/di/sarulabsdi/definition/echo"
	"github.com/tumani1/diexample/di/sarulabsdi/definition/logger"
	"github.com/tumani1/diexample/di/sarulabsdi/internal/definition/postgres"
	"github.com/tumani1/diexample/di/sarulabsdi/internal/domain"
	"github.com/tumani1/diexample/di/sarulabsdi/internal/http/handlers"
)

// DefAutoCompleteHandler name of DI definition
const DefAutoCompleteHandler = "http.handler.autocomplete"

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder) error {
		return builder.Add(container.Definition{
			Name: DefAutoCompleteHandler,
			Tags: []container.Tag{{
				Name: echo.DefHTTPHandlerTag,
			}},
			Build: func(ctx container.Context) (_ interface{}, err error) {
				var log logger.Logger
				if err = ctx.Fill(logger.DefLogger, &log); err != nil {
					return nil, errors.Wrap(err, "can't get logger from container")
				}

				var searchRepo domain.ISearchRepository
				if err = ctx.Fill(postgres.DefSearchPostgresRepo, &searchRepo); err != nil {
					return nil, errors.Wrap(err, "can't get search repo from container")
				}

				var calendarRepo domain.ICalendarRepository
				if err = ctx.Fill(postgres.DefCalendarPostgresRepo, &calendarRepo); err != nil {
					return nil, errors.Wrap(err, "can't get calendar repo from container")
				}

				return handlers.NewSearchHandler(
					log,
					validator.New(),
					searchRepo,
					calendarRepo,
				), nil
			},
		})
	})
}
