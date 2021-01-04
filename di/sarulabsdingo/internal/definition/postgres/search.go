package postgres

import (
	"github.com/tumani1/diexample/di/sarulabsdingo/container"
	"github.com/tumani1/diexample/di/sarulabsdingo/definition/postgres"
	"github.com/tumani1/diexample/di/sarulabsdingo/internal/domain"
	postgresRepo "github.com/tumani1/diexample/di/sarulabsdingo/internal/postgres"
)

// DefSearchPostgresRepo name of DI definition
const DefSearchPostgresRepo = "postgres.repo.search"

// Definition init func.
func init() {
	container.Register(func(builder *container.ProviderObject) error {
		return builder.Add(container.Definition{
			Name: DefSearchPostgresRepo,
			Build: func(db postgres.Postgres) (_ domain.ISearchRepository, err error) {
				return postgresRepo.NewSearchRepository(
					db,
				), nil
			},
			Params: container.Params{
				"0": container.Service(postgres.DefPostgres),
			},
		})
	})
}
