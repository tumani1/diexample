package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/tumani1/diexample/di/sarulabsdingo/definition/logger"
	"github.com/tumani1/diexample/di/sarulabsdingo/definition/postgres"
)

var (
	migrationsPath string

	migrateDatabaseCmd = &cobra.Command{
		Use:   "migrate [command]",
		Short: "Database migrations",
	}

	migrateDatabaseUpCmd = &cobra.Command{
		Use:   "up",
		Short: "Up migrations",
		Args:  cobra.ExactArgs(0),
		RunE:  migrateUpCmdHandler,
	}

	migrateDatabaseDownCmd = &cobra.Command{
		Use:   "down",
		Short: "Down migrations",
		Args:  cobra.ExactArgs(0),
		RunE:  migrateDownCmdHandler,
	}
)

func init() {
	rootCmd.AddCommand(migrateDatabaseCmd)
	migrateDatabaseCmd.AddCommand(migrateDatabaseUpCmd)
	migrateDatabaseCmd.AddCommand(migrateDatabaseDownCmd)

	var err error
	var appPath string

	if appPath, err = os.Getwd(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	migrateDatabaseCmd.PersistentFlags().
		StringVarP(&migrationsPath, "migrationsPath", "m", appPath+"/asstets/migrations", "Path to migrations")

	// Set custom migration table
	migrate.SetTable("migrations")
}

func migrateUpCmdHandler(_ *cobra.Command, _ []string) (err error) {
	migrationsList := &migrate.FileMigrationSource{
		Dir: migrationsPath,
	}

	var db *sql.DB
	if err = diContainer.Fill(postgres.DefPostgres, &db); err != nil {
		return errors.Wrap(err, "can't get postgres from container")
	}

	var log logger.Logger
	if err = diContainer.Fill(logger.DefLogger, &log); err != nil {
		return errors.Wrap(err, "can't get logger from container")
	}

	var n int
	if n, err = migrate.Exec(db, "postgres", migrationsList, migrate.Up); err != nil {
		return errors.Wrap(err, "can't exec migrations")
	}

	log.Info("Applied migrations", zap.Int("count", n))
	return nil
}

func migrateDownCmdHandler(_ *cobra.Command, _ []string) (err error) {
	migrationsList := &migrate.FileMigrationSource{
		Dir: migrationsPath,
	}

	var db *sql.DB
	if err = diContainer.Fill(postgres.DefPostgres, &db); err != nil {
		return
		return errors.Wrap(err, "can't get postgres from container")
	}

	var log logger.Logger
	if err = diContainer.Fill(logger.DefLogger, &log); err != nil {
		return errors.Wrap(err, "can't get logger from container")
	}

	var n int
	if n, err = migrate.ExecMax(db, "postgres", migrationsList, migrate.Down, 1); err != nil {
		return errors.Wrap(err, "can't exec migrations")
	}

	log.Info("Down migrations", zap.Int("count", n))
	return nil
}
