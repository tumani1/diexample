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

// DefSearchHandler name of DI definition
const DefSearchHandler = "http.handler.search"

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Definition{
			Name: DefSearchHandler,
			Tags: []container.Tag{{
				Name: echo.DefHTTPHandlerTag,
			}},
			Build: func(ctx container.Context) (_ interface{}, err error) {
				var log logger.Logger
				if err = ctx.Fill(logger.DefLogger, &log); err != nil {
					return nil, errors.Wrap(err, "can't get logger from container")
				}

				var autoCompleteRepo domain.IAutoCompleteRepository
				if err = ctx.Fill(postgres.DefAutoCompletePostgresRepo, &autoCompleteRepo); err != nil {
					return nil, errors.Wrap(err, "can't get auto complete repo from container")
				}

				return handlers.NewAutoCompleteHandler(
					log,
					validator.New(),
					autoCompleteRepo,
				), nil
			},
		})
	})
}
