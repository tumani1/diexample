package postgres

import (
	"github.com/tumani1/diexample/di/sarulabsdingo/container"
	"github.com/tumani1/diexample/di/sarulabsdingo/definition/postgres"
	"github.com/tumani1/diexample/di/sarulabsdingo/internal/domain"
	postgresRepo "github.com/tumani1/diexample/di/sarulabsdingo/internal/postgres"
)

// DefCalendarPostgresRepo name of DI definition
const DefCalendarPostgresRepo = "postgres.repo.calendar"

// Definition init func.
func init() {
	container.Register(func(builder *container.ProviderObject) error {
		return builder.Add(container.Definition{
			Name: DefCalendarPostgresRepo,
			Build: func(db postgres.Postgres) (_ domain.ICalendarRepository, err error) {
				return postgresRepo.NewCalendarRepository(
					db,
				), nil
			},
			Params: container.Params{
				"0": container.Service(postgres.DefPostgres),
			},
		})
	})
}
