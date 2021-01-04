package config

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/tumani1/diexample/di/sarulabsdi/container"
)

const DefConfig = "config"

type Config = *viper.Viper

func init() {
	container.Register(func(builder *container.Builder) error {
		return builder.Add(container.Definition{
			Name: DefConfig,
			Build: func(ctx container.Context) (interface{}, error) {
				path := os.Getenv("CONFIG_PATH")
				if len(path) == 0 {
					return nil, errors.New("empty config path")
				}

				cfg := viper.New()
				cfg.AutomaticEnv()
				cfg.SetEnvPrefix("ENV")
				cfg.SetEnvKeyReplacer(
					strings.NewReplacer(".", "_"),
				)
				cfg.SetConfigFile(path)
				cfg.SetConfigType("json")

				return cfg, cfg.ReadInConfig()
			},
		})
	})
}
