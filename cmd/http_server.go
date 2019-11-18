package cmd

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"gitlab.com/igor.tumanov1/theboatscom/definition/config"
	"gitlab.com/igor.tumanov1/theboatscom/definition/echo"
	_ "gitlab.com/igor.tumanov1/theboatscom/internal/definition/http/handlers"
)

// http server command.
var httpServerCmd = &cobra.Command{
	Use:   "http-server",
	Short: "Run http server",
	RunE:  httpServerCmdHandler,
}

// Command init function.
func init() {
	rootCmd.AddCommand(httpServerCmd)
}

// Command handler func.
func httpServerCmdHandler(_ *cobra.Command, _ []string) (err error) {
	var conf config.Config
	if err = diContainer.Fill(config.DefConfig, &conf); err != nil {
		return errors.Wrap(err, "can't ger config from container")
	}

	var e echo.Echo
	if err = diContainer.Fill(echo.DefEcho, &e); err != nil {
		return errors.Wrap(err, "can't ger http server from container")
	}

	var s = &http.Server{
		Addr: net.JoinHostPort(
			conf.GetString("http.host"),
			conf.GetString("http.port"),
		),
	}

	go e.Logger.Error(e.StartServer(s))

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return e.Shutdown(ctx)
}
