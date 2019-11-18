package postgres

import (
	"github.com/pkg/errors"

	"gitlab.com/igor.tumanov1/theboatscom/container"
	"gitlab.com/igor.tumanov1/theboatscom/definition/postgres"
	postgresRepo "gitlab.com/igor.tumanov1/theboatscom/internal/postgres"
)

// DefAutoCompletePostgresRepo name of DI definition
const DefAutoCompletePostgresRepo = "postgres.repo.autocomplete"

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Definition{
			Name: DefAutoCompletePostgresRepo,
			Build: func(ctx container.Context) (_ interface{}, err error) {
				var db postgres.Postgres
				if err = ctx.Fill(postgres.DefPostgres, &db); err != nil {
					return nil, errors.Wrap(err, "can't get database from container")
				}

				return postgresRepo.NewAutoCompleteRepository(
					db,
				), nil
			},
		})
	})
}
