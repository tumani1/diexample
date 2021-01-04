package cmd

import (
	"github.com/pkg/errors"
	"github.com/sarulabs/dingo/v4"
	"github.com/spf13/cobra"

	"github.com/tumani1/diexample/di/sarulabsdingo/container"
)

// generate command.
var (
	generatePath string

	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate command",
	}

	generateDefinitionsCmd = &cobra.Command{
		Use:   "definitions",
		Short: "Generate definitions",
		Args:  cobra.ExactArgs(0),
		RunE:  generateCmdCmdHandler,
	}
)

// Command init function.
func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(generateDefinitionsCmd)

	generateCmd.PersistentFlags().
		StringVarP(&generatePath, "generatePath", "g", "", "Path to generation")
}

// Command handler func.
func generateCmdCmdHandler(_ *cobra.Command, _ []string) (err error) {
	provider, err := container.NewProviderObject()
	if err != nil {
		return errors.Wrap(err, "can't create provider")
	}

	err = dingo.GenerateContainer(provider, generatePath)
	if err != nil {
		return errors.Wrap(err, "can't generate container")
	}

	return nil
}
