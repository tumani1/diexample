package cmd

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/tumani1/diexample/di/sarulabsdingo/definition/config"
	"github.com/tumani1/diexample/di/sarulabsdingo/definition/echo"
	appEcho "github.com/tumani1/diexample/di/sarulabsdingo/echo"
	_ "github.com/tumani1/diexample/di/sarulabsdingo/internal/definition/http/handlers"
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
		return errors.Wrap(err, "can't get http server from container")
	}

	for _, def := range diContainer.Definitions() {
		if !strings.Contains(def.Name, echo.DefHTTPHandlerTag) {
			continue
		}

		var c appEcho.Handler
		if err = diContainer.Fill(def.Name, &c); err != nil {
			return errors.Wrap(err, "can't serve handler from container")
		}

		c.Serve(e)
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
