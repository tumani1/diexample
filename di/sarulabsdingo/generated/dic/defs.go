package dic

import (
	"errors"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"

	sql "database/sql"

	echo "github.com/labstack/echo"
	viper "github.com/spf13/viper"
	zap "go.uber.org/zap"

	echo1 "github.com/tumani1/diexample/di/sarulabsdingo/echo"
	domain "github.com/tumani1/diexample/di/sarulabsdingo/internal/domain"
)

func getDiDefs(provider dingo.Provider) []di.Def {
	return []di.Def{
		{
			Name:  "config",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("config")
				if err != nil {
					var eo *viper.Viper
					return eo, err
				}
				b, ok := d.Build.(func() (*viper.Viper, error))
				if !ok {
					var eo *viper.Viper
					return eo, errors.New("could not cast build function to func() (*viper.Viper, error)")
				}
				return b()
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "db.postgres",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("db.postgres")
				if err != nil {
					var eo *sql.DB
					return eo, err
				}
				pi0, err := ctn.SafeGet("config")
				if err != nil {
					var eo *sql.DB
					return eo, err
				}
				p0, ok := pi0.(*viper.Viper)
				if !ok {
					var eo *sql.DB
					return eo, errors.New("could not cast parameter 0 to *viper.Viper")
				}
				b, ok := d.Build.(func(*viper.Viper) (*sql.DB, error))
				if !ok {
					var eo *sql.DB
					return eo, errors.New("could not cast build function to func(*viper.Viper) (*sql.DB, error)")
				}
				return b(p0)
			},
			Close: func(obj interface{}) error {
				d, err := provider.Get("db.postgres")
				if err != nil {
					return err
				}
				c, ok := d.Close.(func(*sql.DB) error)
				if !ok {
					return errors.New("could not cast close function to 'func(*sql.DB) error'")
				}
				o, ok := obj.(*sql.DB)
				if !ok {
					return errors.New("could not cast object to '*sql.DB'")
				}
				return c(o)
			},
		},
		{
			Name:  "echo",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("echo")
				if err != nil {
					var eo *echo.Echo
					return eo, err
				}
				pi0, err := ctn.SafeGet("config")
				if err != nil {
					var eo *echo.Echo
					return eo, err
				}
				p0, ok := pi0.(*viper.Viper)
				if !ok {
					var eo *echo.Echo
					return eo, errors.New("could not cast parameter 0 to *viper.Viper")
				}
				pi1, err := ctn.SafeGet("echo.error_handler")
				if err != nil {
					var eo *echo.Echo
					return eo, err
				}
				p1, ok := pi1.(func(error, echo.Context))
				if !ok {
					var eo *echo.Echo
					return eo, errors.New("could not cast parameter 1 to func(error, echo.Context)")
				}
				b, ok := d.Build.(func(*viper.Viper, func(error, echo.Context)) (*echo.Echo, error))
				if !ok {
					var eo *echo.Echo
					return eo, errors.New("could not cast build function to func(*viper.Viper, func(error, echo.Context)) (*echo.Echo, error)")
				}
				return b(p0, p1)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "echo.error_handler",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("echo.error_handler")
				if err != nil {
					var eo func(error, echo.Context)
					return eo, err
				}
				b, ok := d.Build.(func() (func(error, echo.Context), error))
				if !ok {
					var eo func(error, echo.Context)
					return eo, errors.New("could not cast build function to func() (func(error, echo.Context), error)")
				}
				return b()
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "http.handler.autocomplete",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("http.handler.autocomplete")
				if err != nil {
					var eo echo1.Handler
					return eo, err
				}
				pi0, err := ctn.SafeGet("logger")
				if err != nil {
					var eo echo1.Handler
					return eo, err
				}
				p0, ok := pi0.(*zap.Logger)
				if !ok {
					var eo echo1.Handler
					return eo, errors.New("could not cast parameter 0 to *zap.Logger")
				}
				pi1, err := ctn.SafeGet("postgres.repo.search")
				if err != nil {
					var eo echo1.Handler
					return eo, err
				}
				p1, ok := pi1.(domain.ISearchRepository)
				if !ok {
					var eo echo1.Handler
					return eo, errors.New("could not cast parameter 1 to domain.ISearchRepository")
				}
				pi2, err := ctn.SafeGet("postgres.repo.calendar")
				if err != nil {
					var eo echo1.Handler
					return eo, err
				}
				p2, ok := pi2.(domain.ICalendarRepository)
				if !ok {
					var eo echo1.Handler
					return eo, errors.New("could not cast parameter 2 to domain.ICalendarRepository")
				}
				b, ok := d.Build.(func(*zap.Logger, domain.ISearchRepository, domain.ICalendarRepository) (echo1.Handler, error))
				if !ok {
					var eo echo1.Handler
					return eo, errors.New("could not cast build function to func(*zap.Logger, domain.ISearchRepository, domain.ICalendarRepository) (echo1.Handler, error)")
				}
				return b(p0, p1, p2)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "http.handler.search",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("http.handler.search")
				if err != nil {
					var eo echo1.Handler
					return eo, err
				}
				pi0, err := ctn.SafeGet("logger")
				if err != nil {
					var eo echo1.Handler
					return eo, err
				}
				p0, ok := pi0.(*zap.Logger)
				if !ok {
					var eo echo1.Handler
					return eo, errors.New("could not cast parameter 0 to *zap.Logger")
				}
				pi1, err := ctn.SafeGet("postgres.repo.autocomplete")
				if err != nil {
					var eo echo1.Handler
					return eo, err
				}
				p1, ok := pi1.(domain.IAutoCompleteRepository)
				if !ok {
					var eo echo1.Handler
					return eo, errors.New("could not cast parameter 1 to domain.IAutoCompleteRepository")
				}
				b, ok := d.Build.(func(*zap.Logger, domain.IAutoCompleteRepository) (echo1.Handler, error))
				if !ok {
					var eo echo1.Handler
					return eo, errors.New("could not cast build function to func(*zap.Logger, domain.IAutoCompleteRepository) (echo1.Handler, error)")
				}
				return b(p0, p1)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "logger",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("logger")
				if err != nil {
					var eo *zap.Logger
					return eo, err
				}
				pi0, err := ctn.SafeGet("config")
				if err != nil {
					var eo *zap.Logger
					return eo, err
				}
				p0, ok := pi0.(*viper.Viper)
				if !ok {
					var eo *zap.Logger
					return eo, errors.New("could not cast parameter 0 to *viper.Viper")
				}
				b, ok := d.Build.(func(*viper.Viper) (*zap.Logger, error))
				if !ok {
					var eo *zap.Logger
					return eo, errors.New("could not cast build function to func(*viper.Viper) (*zap.Logger, error)")
				}
				return b(p0)
			},
			Close: func(obj interface{}) error {
				d, err := provider.Get("logger")
				if err != nil {
					return err
				}
				c, ok := d.Close.(func(*zap.Logger) error)
				if !ok {
					return errors.New("could not cast close function to 'func(*zap.Logger) error'")
				}
				o, ok := obj.(*zap.Logger)
				if !ok {
					return errors.New("could not cast object to '*zap.Logger'")
				}
				return c(o)
			},
		},
		{
			Name:  "postgres.repo.autocomplete",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("postgres.repo.autocomplete")
				if err != nil {
					var eo domain.IAutoCompleteRepository
					return eo, err
				}
				pi0, err := ctn.SafeGet("db.postgres")
				if err != nil {
					var eo domain.IAutoCompleteRepository
					return eo, err
				}
				p0, ok := pi0.(*sql.DB)
				if !ok {
					var eo domain.IAutoCompleteRepository
					return eo, errors.New("could not cast parameter 0 to *sql.DB")
				}
				b, ok := d.Build.(func(*sql.DB) (domain.IAutoCompleteRepository, error))
				if !ok {
					var eo domain.IAutoCompleteRepository
					return eo, errors.New("could not cast build function to func(*sql.DB) (domain.IAutoCompleteRepository, error)")
				}
				return b(p0)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "postgres.repo.calendar",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("postgres.repo.calendar")
				if err != nil {
					var eo domain.ICalendarRepository
					return eo, err
				}
				pi0, err := ctn.SafeGet("db.postgres")
				if err != nil {
					var eo domain.ICalendarRepository
					return eo, err
				}
				p0, ok := pi0.(*sql.DB)
				if !ok {
					var eo domain.ICalendarRepository
					return eo, errors.New("could not cast parameter 0 to *sql.DB")
				}
				b, ok := d.Build.(func(*sql.DB) (domain.ICalendarRepository, error))
				if !ok {
					var eo domain.ICalendarRepository
					return eo, errors.New("could not cast build function to func(*sql.DB) (domain.ICalendarRepository, error)")
				}
				return b(p0)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
		{
			Name:  "postgres.repo.search",
			Scope: "",
			Build: func(ctn di.Container) (interface{}, error) {
				d, err := provider.Get("postgres.repo.search")
				if err != nil {
					var eo domain.ISearchRepository
					return eo, err
				}
				pi0, err := ctn.SafeGet("db.postgres")
				if err != nil {
					var eo domain.ISearchRepository
					return eo, err
				}
				p0, ok := pi0.(*sql.DB)
				if !ok {
					var eo domain.ISearchRepository
					return eo, errors.New("could not cast parameter 0 to *sql.DB")
				}
				b, ok := d.Build.(func(*sql.DB) (domain.ISearchRepository, error))
				if !ok {
					var eo domain.ISearchRepository
					return eo, errors.New("could not cast build function to func(*sql.DB) (domain.ISearchRepository, error)")
				}
				return b(p0)
			},
			Close: func(obj interface{}) error {
				return nil
			},
		},
	}
}
