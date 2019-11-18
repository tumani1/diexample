package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"

	"gitlab.com/igor.tumanov1/theboatscom/container"
)

const DefConfig = "config"

type Config = *viper.Viper

func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		var ok bool
		if _, ok = params["config_path"]; !ok {
			return errors.New("can't get required parameter config path")
		}

		var path string
		if path, ok = params["config_path"].(string); !ok {
			return errors.New(`parameter "config_path" should be string`)
		}

		return builder.Add(container.Definition{
			Name: DefConfig,
			Build: func(ctx container.Context) (interface{}, error) {
				var cfg = viper.New()

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
