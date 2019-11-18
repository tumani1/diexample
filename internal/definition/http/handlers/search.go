package handlers

import (
	"github.com/pkg/errors"
	"gopkg.in/go-playground/validator.v9"

	"gitlab.com/igor.tumanov1/theboatscom/container"
	"gitlab.com/igor.tumanov1/theboatscom/definition/echo"
	"gitlab.com/igor.tumanov1/theboatscom/definition/logger"
	"gitlab.com/igor.tumanov1/theboatscom/internal/definition/postgres"
	"gitlab.com/igor.tumanov1/theboatscom/internal/domain"
	"gitlab.com/igor.tumanov1/theboatscom/internal/http/handlers"
)

// DefAutoCompleteHandler name of DI definition
const DefAutoCompleteHandler = "http.handler.autocomplete"

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
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
