package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/tumani1/diexample/di/sarulabsdi/container"
	"github.com/tumani1/diexample/di/sarulabsdi/definition/config"
)

const DefPostgres = "db.postgres"

type Postgres = *sql.DB

func init() {
	container.Register(func(builder *container.Builder) error {
		return builder.Add(container.Definition{
			Name: DefPostgres,
			Build: func(ctx container.Context) (_ interface{}, err error) {
				var cfg config.Config
				if err = ctx.Fill(config.DefConfig, &cfg); err != nil {
					return nil, errors.Wrap(err, "can't get config from container")
				}

				var db *sql.DB
				if db, err = sql.Open("postgres", cfg.GetString("postgres.server")); err != nil {
					return nil, errors.Wrap(err, "can't open db connection")
				}

				// setup max open connections
				db.SetMaxOpenConns(cfg.GetInt("postgres.max_open_conns"))

				// setup max idle connections
				db.SetMaxIdleConns(cfg.GetInt("postgres.max_idle_conns"))

				// setup max life time connection
				db.SetConnMaxLifetime(cfg.GetDuration("postgres.max_life_time_conn"))

				if err = db.Ping(); err != nil {
					return nil, errors.Wrap(err, "Error ping postgres")
				}

				return db, nil
			},
			Close: func(obj interface{}) (err error) {
				if err = obj.(*sql.DB).Close(); err != nil {
					return errors.Wrap(err, "can't close postgres connection")
				}

				return nil
			},
		})
	})
}
