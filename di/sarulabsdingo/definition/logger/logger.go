package logger

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/tumani1/diexample/di/sarulabsdingo/container"
	"github.com/tumani1/diexample/di/sarulabsdingo/definition/config"
)

const DefLogger = "logger"

type (
	Logger = *zap.Logger

	loggerConf struct {
		Cores []struct {
			Addr     string
			Host     string
			Level    string
			Encoding string
		}
		Caller      bool
		Stacktrace  string
		Development bool
	}
)

func init() {
	container.Register(func(builder *container.ProviderObject) error {
		return builder.Add(container.Definition{
			Name: DefLogger,
			Build: func(cfg config.Config) (_ *zap.Logger, err error) {
				var conf loggerConf
				if err = cfg.UnmarshalKey("logger", &conf); err != nil {
					return nil, err
				}

				var cores = make([]zapcore.Core, 0, 2)
				for _, logger := range conf.Cores {
					var core zapcore.Core
					switch logger.Encoding {
					case "console", "json":
						var eConf zapcore.EncoderConfig
						if conf.Development {
							eConf = zap.NewDevelopmentEncoderConfig()
						} else {
							eConf = zap.NewProductionEncoderConfig()
						}

						var level zap.AtomicLevel
						if len(logger.Level) > 0 {
							if err = level.UnmarshalText([]byte(logger.Level)); err != nil {
								return nil, err
							}
						}

						var enc = zapcore.NewConsoleEncoder(eConf)
						if logger.Encoding == "json" {
							enc = zapcore.NewJSONEncoder(eConf)
						}

						core = zapcore.NewCore(enc, zapcore.AddSync(os.Stdout), level)
					default:
						return nil, fmt.Errorf("unknown encoding %s", logger.Encoding)
					}

					cores = append(cores, core)
				}

				if len(cores) == 0 {
					cores = append(cores, zapcore.NewNopCore())
				}

				var options = make([]zap.Option, 0, 8)
				if conf.Caller {
					options = append(options, zap.AddCaller())
				}

				if conf.Development {
					options = append(options, zap.Development())
				}

				var level zap.AtomicLevel
				if len(conf.Stacktrace) > 0 {
					if err = level.UnmarshalText([]byte(conf.Stacktrace)); err != nil {
						return nil, err
					}

					options = append(options, zap.AddStacktrace(level))
				}

				return zap.New(zapcore.NewTee(cores...), options...), nil
			},
			Params: container.Params{
				"0": container.Service(config.DefConfig),
			},
			Close: func(obj *zap.Logger) error {
				if err := obj.Sync(); err != nil {
					return errors.Wrap(err, "error sync logger")
				}

				return nil
			},
		})
	})
}
