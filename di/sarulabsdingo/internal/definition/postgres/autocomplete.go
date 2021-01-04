package postgres

import (
	"github.com/tumani1/diexample/di/sarulabsdingo/container"
	"github.com/tumani1/diexample/di/sarulabsdingo/definition/postgres"
	"github.com/tumani1/diexample/di/sarulabsdingo/internal/domain"
	postgresRepo "github.com/tumani1/diexample/di/sarulabsdingo/internal/postgres"
)

// DefAutoCompletePostgresRepo name of DI definition
const DefAutoCompletePostgresRepo = "postgres.repo.autocomplete"

// Definition init func.
func init() {
	container.Register(func(builder *container.ProviderObject) error {
		return builder.Add(container.Definition{
			Name: DefAutoCompletePostgresRepo,
			Build: func(db postgres.Postgres) (_ domain.IAutoCompleteRepository, err error) {
				return postgresRepo.NewAutoCompleteRepository(
					db,
				), nil
			},
			Params: container.Params{
				"0": container.Service(postgres.DefPostgres),
			},
		})
	})
}
