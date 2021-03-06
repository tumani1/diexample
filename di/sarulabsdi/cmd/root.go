package cmd

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/tumani1/diexample/di/sarulabsdi/container"
)

var (
	// Config path.
	configPath string

	// DI Container.
	diContainer container.Context

	// Root command.
	rootCmd = &cobra.Command{
		Use:           "app [command]",
		Long:          "",
		Short:         "",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {
			diContainer, err = container.Instance()
			if err != nil {
				return errors.Wrap(err, "can't get container")
			}

			return nil
		},
	}
)

func Execute() (err error) {
	var appPath string
	if appPath, err = os.Getwd(); err != nil {
		return errors.Wrap(err, "can't get working dir")
	}

	// Application config path
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", appPath+"/configs/config.json", "config file")

	// Run
	err = rootCmd.Execute()

	// Delete context
	if diContainer != nil {
		_ = diContainer.Delete()
	}

	return err
}
