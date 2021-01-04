package postgres

import (
	"github.com/pkg/errors"

	"github.com/tumani1/diexample/di/sarulabsdi/container"
	"github.com/tumani1/diexample/di/sarulabsdi/definition/postgres"
	postgresRepo "github.com/tumani1/diexample/di/sarulabsdi/internal/postgres"
)

// DefCalendarPostgresRepo name of DI definition
const DefCalendarPostgresRepo = "postgres.repo.calendar"

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder) error {
		return builder.Add(container.Definition{
			Name: DefCalendarPostgresRepo,
			Build: func(ctx container.Context) (_ interface{}, err error) {
				var db postgres.Postgres
				if err = ctx.Fill(postgres.DefPostgres, &db); err != nil {
					return nil, errors.Wrap(err, "can't get database from container")
				}

				return postgresRepo.NewCalendarRepository(
					db,
				), nil
			},
		})
	})
}
