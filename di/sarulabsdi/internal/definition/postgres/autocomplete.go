package postgres

import (
	"github.com/pkg/errors"

	"github.com/tumani1/diexample/di/sarulabsdi/container"
	"github.com/tumani1/diexample/di/sarulabsdi/definition/postgres"
	postgresRepo "github.com/tumani1/diexample/di/sarulabsdi/internal/postgres"
)

// DefAutoCompletePostgresRepo name of DI definition
const DefAutoCompletePostgresRepo = "postgres.repo.autocomplete"

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder) error {
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
