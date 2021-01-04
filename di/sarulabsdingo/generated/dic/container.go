package dic

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"

	providerPkg "github.com/tumani1/diexample/di/sarulabsdingo/container"

	sql "database/sql"

	echo "github.com/labstack/echo"
	viper "github.com/spf13/viper"
	zap "go.uber.org/zap"

	echo1 "github.com/tumani1/diexample/di/sarulabsdingo/echo"
	domain "github.com/tumani1/diexample/di/sarulabsdingo/internal/domain"
)

// C retrieves a Container from an interface.
// The function panics if the Container can not be retrieved.
//
// The interface can be :
// - a *Container
// - an *http.Request containing a *Container in its context.Context
//   for the dingo.ContainerKey("dingo") key.
//
// The function can be changed to match the needs of your application.
var C = func(i interface{}) *Container {
	if c, ok := i.(*Container); ok {
		return c
	}
	r, ok := i.(*http.Request)
	if !ok {
		panic("could not get the container with dic.C()")
	}
	c, ok := r.Context().Value(dingo.ContainerKey("dingo")).(*Container)
	if !ok {
		panic("could not get the container from the given *http.Request in dic.C()")
	}
	return c
}

type builder struct {
	builder *di.Builder
}

// NewBuilder creates a builder that can create a Container.
// You should you NewContainer to create the container directly.
// Using NewBuilder allows you to redefine some di services though.
// This could be used for testing.
// But this behaviour is not safe, so be sure to know what you are doing.
func NewBuilder(scopes ...string) (*builder, error) {
	if len(scopes) == 0 {
		scopes = []string{di.App, di.Request, di.SubRequest}
	}
	b, err := di.NewBuilder(scopes...)
	if err != nil {
		return nil, fmt.Errorf("could not create di.Builder: %v", err)
	}
	provider := &providerPkg.ProviderObject{}
	if err := provider.Load(); err != nil {
		return nil, fmt.Errorf("could not load definitions with the Provider (ProviderObject from gitlab.com/igor.tumanov1/theboatscom/di/sarulabsdingo/container): %v", err)
	}
	for _, d := range getDiDefs(provider) {
		if err := b.Add(d); err != nil {
			return nil, fmt.Errorf("could not add di.Def in di.Builder: %v", err)
		}
	}
	return &builder{builder: b}, nil
}

// Add adds one or more definitions in the Builder.
// It returns an error if a definition can not be added.
func (b *builder) Add(defs ...di.Def) error {
	return b.builder.Add(defs...)
}

// Set is a shortcut to add a definition for an already built object.
func (b *builder) Set(name string, obj interface{}) error {
	return b.builder.Set(name, obj)
}

// Build creates a Container in the most generic scope.
func (b *builder) Build() *Container {
	return &Container{ctn: b.builder.Build()}
}

// NewContainer creates a new Container.
// If no scope is provided, di.App, di.Request and di.SubRequest are used.
// The returned Container has the most generic scope (di.App).
// The SubContainer() method should be called to get a Container in a more specific scope.
func NewContainer(scopes ...string) (*Container, error) {
	b, err := NewBuilder(scopes...)
	if err != nil {
		return nil, err
	}
	return b.Build(), nil
}

// Container represents a generated dependency injection container.
// It is a wrapper around a di.Container.
//
// A Container has a scope and may have a parent in a more generic scope
// and children in a more specific scope.
// Objects can be retrieved from the Container.
// If the requested object does not already exist in the Container,
// it is built thanks to the object definition.
// The following attempts to get this object will return the same object.
type Container struct {
	ctn di.Container
}

func (c *Container) Definitions() map[string]di.Def {
	return c.ctn.Definitions()
}

// Scope returns the Container scope.
func (c *Container) Scope() string {
	return c.ctn.Scope()
}

// Scopes returns the list of available scopes.
func (c *Container) Scopes() []string {
	return c.ctn.Scopes()
}

// ParentScopes returns the list of scopes wider than the Container scope.
func (c *Container) ParentScopes() []string {
	return c.ctn.ParentScopes()
}

// SubScopes returns the list of scopes that are more specific than the Container scope.
func (c *Container) SubScopes() []string {
	return c.ctn.SubScopes()
}

// Parent returns the parent Container.
func (c *Container) Parent() *Container {
	if p := c.ctn.Parent(); p != nil {
		return &Container{ctn: p}
	}
	return nil
}

// SubContainer creates a new Container in the next sub-scope
// that will have this Container as parent.
func (c *Container) SubContainer() (*Container, error) {
	sub, err := c.ctn.SubContainer()
	if err != nil {
		return nil, err
	}
	return &Container{ctn: sub}, nil
}

// SafeGet retrieves an object from the Container.
// The object has to belong to this scope or a more generic one.
// If the object does not already exist, it is created and saved in the Container.
// If the object can not be created, it returns an error.
func (c *Container) SafeGet(name string) (interface{}, error) {
	return c.ctn.SafeGet(name)
}

// Get is similar to SafeGet but it does not return the error.
// Instead it panics.
func (c *Container) Get(name string) interface{} {
	return c.ctn.Get(name)
}

// Fill is similar to SafeGet but it does not return the object.
// Instead it fills the provided object with the value returned by SafeGet.
// The provided object must be a pointer to the value returned by SafeGet.
func (c *Container) Fill(name string, dst interface{}) error {
	return c.ctn.Fill(name, dst)
}

// UnscopedSafeGet retrieves an object from the Container, like SafeGet.
// The difference is that the object can be retrieved
// even if it belongs to a more specific scope.
// To do so, UnscopedSafeGet creates a sub-container.
// When the created object is no longer needed,
// it is important to use the Clean method to delete this sub-container.
func (c *Container) UnscopedSafeGet(name string) (interface{}, error) {
	return c.ctn.UnscopedSafeGet(name)
}

// UnscopedGet is similar to UnscopedSafeGet but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGet(name string) interface{} {
	return c.ctn.UnscopedGet(name)
}

// UnscopedFill is similar to UnscopedSafeGet but copies the object in dst instead of returning it.
func (c *Container) UnscopedFill(name string, dst interface{}) error {
	return c.ctn.UnscopedFill(name, dst)
}

// Clean deletes the sub-container created by UnscopedSafeGet, UnscopedGet or UnscopedFill.
func (c *Container) Clean() error {
	return c.ctn.Clean()
}

// DeleteWithSubContainers takes all the objects saved in this Container
// and calls the Close function of their Definition on them.
// It will also call DeleteWithSubContainers on each child and remove its reference in the parent Container.
// After deletion, the Container can no longer be used.
// The sub-containers are deleted even if they are still used in other goroutines.
// It can cause errors. You may want to use the Delete method instead.
func (c *Container) DeleteWithSubContainers() error {
	return c.ctn.DeleteWithSubContainers()
}

// Delete works like DeleteWithSubContainers if the Container does not have any child.
// But if the Container has sub-containers, it will not be deleted right away.
// The deletion only occurs when all the sub-containers have been deleted manually.
// So you have to call Delete or DeleteWithSubContainers on all the sub-containers.
func (c *Container) Delete() error {
	return c.ctn.Delete()
}

// IsClosed returns true if the Container has been deleted.
func (c *Container) IsClosed() bool {
	return c.ctn.IsClosed()
}

// SafeGetConfig works like SafeGet but only for Config.
// It does not return an interface but a *viper.Viper.
func (c *Container) SafeGetConfig() (*viper.Viper, error) {
	i, err := c.ctn.SafeGet("config")
	if err != nil {
		var eo *viper.Viper
		return eo, err
	}
	o, ok := i.(*viper.Viper)
	if !ok {
		return o, errors.New("could get 'config' because the object could not be cast to *viper.Viper")
	}
	return o, nil
}

// GetConfig is similar to SafeGetConfig but it does not return the error.
// Instead it panics.
func (c *Container) GetConfig() *viper.Viper {
	o, err := c.SafeGetConfig()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetConfig works like UnscopedSafeGet but only for Config.
// It does not return an interface but a *viper.Viper.
func (c *Container) UnscopedSafeGetConfig() (*viper.Viper, error) {
	i, err := c.ctn.UnscopedSafeGet("config")
	if err != nil {
		var eo *viper.Viper
		return eo, err
	}
	o, ok := i.(*viper.Viper)
	if !ok {
		return o, errors.New("could get 'config' because the object could not be cast to *viper.Viper")
	}
	return o, nil
}

// UnscopedGetConfig is similar to UnscopedSafeGetConfig but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetConfig() *viper.Viper {
	o, err := c.UnscopedSafeGetConfig()
	if err != nil {
		panic(err)
	}
	return o
}

// Config is similar to GetConfig.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetConfig method.
// If the container can not be retrieved, it panics.
func Config(i interface{}) *viper.Viper {
	return C(i).GetConfig()
}

// SafeGetDbPostgres works like SafeGet but only for DbPostgres.
// It does not return an interface but a *sql.DB.
func (c *Container) SafeGetDbPostgres() (*sql.DB, error) {
	i, err := c.ctn.SafeGet("db.postgres")
	if err != nil {
		var eo *sql.DB
		return eo, err
	}
	o, ok := i.(*sql.DB)
	if !ok {
		return o, errors.New("could get 'db.postgres' because the object could not be cast to *sql.DB")
	}
	return o, nil
}

// GetDbPostgres is similar to SafeGetDbPostgres but it does not return the error.
// Instead it panics.
func (c *Container) GetDbPostgres() *sql.DB {
	o, err := c.SafeGetDbPostgres()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetDbPostgres works like UnscopedSafeGet but only for DbPostgres.
// It does not return an interface but a *sql.DB.
func (c *Container) UnscopedSafeGetDbPostgres() (*sql.DB, error) {
	i, err := c.ctn.UnscopedSafeGet("db.postgres")
	if err != nil {
		var eo *sql.DB
		return eo, err
	}
	o, ok := i.(*sql.DB)
	if !ok {
		return o, errors.New("could get 'db.postgres' because the object could not be cast to *sql.DB")
	}
	return o, nil
}

// UnscopedGetDbPostgres is similar to UnscopedSafeGetDbPostgres but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetDbPostgres() *sql.DB {
	o, err := c.UnscopedSafeGetDbPostgres()
	if err != nil {
		panic(err)
	}
	return o
}

// DbPostgres is similar to GetDbPostgres.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetDbPostgres method.
// If the container can not be retrieved, it panics.
func DbPostgres(i interface{}) *sql.DB {
	return C(i).GetDbPostgres()
}

// SafeGetEcho works like SafeGet but only for Echo.
// It does not return an interface but a *echo.Echo.
func (c *Container) SafeGetEcho() (*echo.Echo, error) {
	i, err := c.ctn.SafeGet("echo")
	if err != nil {
		var eo *echo.Echo
		return eo, err
	}
	o, ok := i.(*echo.Echo)
	if !ok {
		return o, errors.New("could get 'echo' because the object could not be cast to *echo.Echo")
	}
	return o, nil
}

// GetEcho is similar to SafeGetEcho but it does not return the error.
// Instead it panics.
func (c *Container) GetEcho() *echo.Echo {
	o, err := c.SafeGetEcho()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetEcho works like UnscopedSafeGet but only for Echo.
// It does not return an interface but a *echo.Echo.
func (c *Container) UnscopedSafeGetEcho() (*echo.Echo, error) {
	i, err := c.ctn.UnscopedSafeGet("echo")
	if err != nil {
		var eo *echo.Echo
		return eo, err
	}
	o, ok := i.(*echo.Echo)
	if !ok {
		return o, errors.New("could get 'echo' because the object could not be cast to *echo.Echo")
	}
	return o, nil
}

// UnscopedGetEcho is similar to UnscopedSafeGetEcho but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetEcho() *echo.Echo {
	o, err := c.UnscopedSafeGetEcho()
	if err != nil {
		panic(err)
	}
	return o
}

// Echo is similar to GetEcho.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetEcho method.
// If the container can not be retrieved, it panics.
func Echo(i interface{}) *echo.Echo {
	return C(i).GetEcho()
}

// SafeGetEchoErrorHandler works like SafeGet but only for EchoErrorHandler.
// It does not return an interface but a func(error, echo.Context).
func (c *Container) SafeGetEchoErrorHandler() (func(error, echo.Context), error) {
	i, err := c.ctn.SafeGet("echo.error_handler")
	if err != nil {
		var eo func(error, echo.Context)
		return eo, err
	}
	o, ok := i.(func(error, echo.Context))
	if !ok {
		return o, errors.New("could get 'echo.error_handler' because the object could not be cast to func(error, echo.Context)")
	}
	return o, nil
}

// GetEchoErrorHandler is similar to SafeGetEchoErrorHandler but it does not return the error.
// Instead it panics.
func (c *Container) GetEchoErrorHandler() func(error, echo.Context) {
	o, err := c.SafeGetEchoErrorHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetEchoErrorHandler works like UnscopedSafeGet but only for EchoErrorHandler.
// It does not return an interface but a func(error, echo.Context).
func (c *Container) UnscopedSafeGetEchoErrorHandler() (func(error, echo.Context), error) {
	i, err := c.ctn.UnscopedSafeGet("echo.error_handler")
	if err != nil {
		var eo func(error, echo.Context)
		return eo, err
	}
	o, ok := i.(func(error, echo.Context))
	if !ok {
		return o, errors.New("could get 'echo.error_handler' because the object could not be cast to func(error, echo.Context)")
	}
	return o, nil
}

// UnscopedGetEchoErrorHandler is similar to UnscopedSafeGetEchoErrorHandler but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetEchoErrorHandler() func(error, echo.Context) {
	o, err := c.UnscopedSafeGetEchoErrorHandler()
	if err != nil {
		panic(err)
	}
	return o
}

// EchoErrorHandler is similar to GetEchoErrorHandler.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetEchoErrorHandler method.
// If the container can not be retrieved, it panics.
func EchoErrorHandler(i interface{}) func(error, echo.Context) {
	return C(i).GetEchoErrorHandler()
}

// SafeGetHttpHandlerAutocomplete works like SafeGet but only for HttpHandlerAutocomplete.
// It does not return an interface but a echo1.Handler.
func (c *Container) SafeGetHttpHandlerAutocomplete() (echo1.Handler, error) {
	i, err := c.ctn.SafeGet("http.handler.autocomplete")
	if err != nil {
		var eo echo1.Handler
		return eo, err
	}
	o, ok := i.(echo1.Handler)
	if !ok {
		return o, errors.New("could get 'http.handler.autocomplete' because the object could not be cast to echo1.Handler")
	}
	return o, nil
}

// GetHttpHandlerAutocomplete is similar to SafeGetHttpHandlerAutocomplete but it does not return the error.
// Instead it panics.
func (c *Container) GetHttpHandlerAutocomplete() echo1.Handler {
	o, err := c.SafeGetHttpHandlerAutocomplete()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetHttpHandlerAutocomplete works like UnscopedSafeGet but only for HttpHandlerAutocomplete.
// It does not return an interface but a echo1.Handler.
func (c *Container) UnscopedSafeGetHttpHandlerAutocomplete() (echo1.Handler, error) {
	i, err := c.ctn.UnscopedSafeGet("http.handler.autocomplete")
	if err != nil {
		var eo echo1.Handler
		return eo, err
	}
	o, ok := i.(echo1.Handler)
	if !ok {
		return o, errors.New("could get 'http.handler.autocomplete' because the object could not be cast to echo1.Handler")
	}
	return o, nil
}

// UnscopedGetHttpHandlerAutocomplete is similar to UnscopedSafeGetHttpHandlerAutocomplete but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetHttpHandlerAutocomplete() echo1.Handler {
	o, err := c.UnscopedSafeGetHttpHandlerAutocomplete()
	if err != nil {
		panic(err)
	}
	return o
}

// HttpHandlerAutocomplete is similar to GetHttpHandlerAutocomplete.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetHttpHandlerAutocomplete method.
// If the container can not be retrieved, it panics.
func HttpHandlerAutocomplete(i interface{}) echo1.Handler {
	return C(i).GetHttpHandlerAutocomplete()
}

// SafeGetHttpHandlerSearch works like SafeGet but only for HttpHandlerSearch.
// It does not return an interface but a echo1.Handler.
func (c *Container) SafeGetHttpHandlerSearch() (echo1.Handler, error) {
	i, err := c.ctn.SafeGet("http.handler.search")
	if err != nil {
		var eo echo1.Handler
		return eo, err
	}
	o, ok := i.(echo1.Handler)
	if !ok {
		return o, errors.New("could get 'http.handler.search' because the object could not be cast to echo1.Handler")
	}
	return o, nil
}

// GetHttpHandlerSearch is similar to SafeGetHttpHandlerSearch but it does not return the error.
// Instead it panics.
func (c *Container) GetHttpHandlerSearch() echo1.Handler {
	o, err := c.SafeGetHttpHandlerSearch()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetHttpHandlerSearch works like UnscopedSafeGet but only for HttpHandlerSearch.
// It does not return an interface but a echo1.Handler.
func (c *Container) UnscopedSafeGetHttpHandlerSearch() (echo1.Handler, error) {
	i, err := c.ctn.UnscopedSafeGet("http.handler.search")
	if err != nil {
		var eo echo1.Handler
		return eo, err
	}
	o, ok := i.(echo1.Handler)
	if !ok {
		return o, errors.New("could get 'http.handler.search' because the object could not be cast to echo1.Handler")
	}
	return o, nil
}

// UnscopedGetHttpHandlerSearch is similar to UnscopedSafeGetHttpHandlerSearch but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetHttpHandlerSearch() echo1.Handler {
	o, err := c.UnscopedSafeGetHttpHandlerSearch()
	if err != nil {
		panic(err)
	}
	return o
}

// HttpHandlerSearch is similar to GetHttpHandlerSearch.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetHttpHandlerSearch method.
// If the container can not be retrieved, it panics.
func HttpHandlerSearch(i interface{}) echo1.Handler {
	return C(i).GetHttpHandlerSearch()
}

// SafeGetLogger works like SafeGet but only for Logger.
// It does not return an interface but a *zap.Logger.
func (c *Container) SafeGetLogger() (*zap.Logger, error) {
	i, err := c.ctn.SafeGet("logger")
	if err != nil {
		var eo *zap.Logger
		return eo, err
	}
	o, ok := i.(*zap.Logger)
	if !ok {
		return o, errors.New("could get 'logger' because the object could not be cast to *zap.Logger")
	}
	return o, nil
}

// GetLogger is similar to SafeGetLogger but it does not return the error.
// Instead it panics.
func (c *Container) GetLogger() *zap.Logger {
	o, err := c.SafeGetLogger()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetLogger works like UnscopedSafeGet but only for Logger.
// It does not return an interface but a *zap.Logger.
func (c *Container) UnscopedSafeGetLogger() (*zap.Logger, error) {
	i, err := c.ctn.UnscopedSafeGet("logger")
	if err != nil {
		var eo *zap.Logger
		return eo, err
	}
	o, ok := i.(*zap.Logger)
	if !ok {
		return o, errors.New("could get 'logger' because the object could not be cast to *zap.Logger")
	}
	return o, nil
}

// UnscopedGetLogger is similar to UnscopedSafeGetLogger but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetLogger() *zap.Logger {
	o, err := c.UnscopedSafeGetLogger()
	if err != nil {
		panic(err)
	}
	return o
}

// Logger is similar to GetLogger.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetLogger method.
// If the container can not be retrieved, it panics.
func Logger(i interface{}) *zap.Logger {
	return C(i).GetLogger()
}

// SafeGetPostgresRepoAutocomplete works like SafeGet but only for PostgresRepoAutocomplete.
// It does not return an interface but a domain.IAutoCompleteRepository.
func (c *Container) SafeGetPostgresRepoAutocomplete() (domain.IAutoCompleteRepository, error) {
	i, err := c.ctn.SafeGet("postgres.repo.autocomplete")
	if err != nil {
		var eo domain.IAutoCompleteRepository
		return eo, err
	}
	o, ok := i.(domain.IAutoCompleteRepository)
	if !ok {
		return o, errors.New("could get 'postgres.repo.autocomplete' because the object could not be cast to domain.IAutoCompleteRepository")
	}
	return o, nil
}

// GetPostgresRepoAutocomplete is similar to SafeGetPostgresRepoAutocomplete but it does not return the error.
// Instead it panics.
func (c *Container) GetPostgresRepoAutocomplete() domain.IAutoCompleteRepository {
	o, err := c.SafeGetPostgresRepoAutocomplete()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetPostgresRepoAutocomplete works like UnscopedSafeGet but only for PostgresRepoAutocomplete.
// It does not return an interface but a domain.IAutoCompleteRepository.
func (c *Container) UnscopedSafeGetPostgresRepoAutocomplete() (domain.IAutoCompleteRepository, error) {
	i, err := c.ctn.UnscopedSafeGet("postgres.repo.autocomplete")
	if err != nil {
		var eo domain.IAutoCompleteRepository
		return eo, err
	}
	o, ok := i.(domain.IAutoCompleteRepository)
	if !ok {
		return o, errors.New("could get 'postgres.repo.autocomplete' because the object could not be cast to domain.IAutoCompleteRepository")
	}
	return o, nil
}

// UnscopedGetPostgresRepoAutocomplete is similar to UnscopedSafeGetPostgresRepoAutocomplete but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetPostgresRepoAutocomplete() domain.IAutoCompleteRepository {
	o, err := c.UnscopedSafeGetPostgresRepoAutocomplete()
	if err != nil {
		panic(err)
	}
	return o
}

// PostgresRepoAutocomplete is similar to GetPostgresRepoAutocomplete.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetPostgresRepoAutocomplete method.
// If the container can not be retrieved, it panics.
func PostgresRepoAutocomplete(i interface{}) domain.IAutoCompleteRepository {
	return C(i).GetPostgresRepoAutocomplete()
}

// SafeGetPostgresRepoCalendar works like SafeGet but only for PostgresRepoCalendar.
// It does not return an interface but a domain.ICalendarRepository.
func (c *Container) SafeGetPostgresRepoCalendar() (domain.ICalendarRepository, error) {
	i, err := c.ctn.SafeGet("postgres.repo.calendar")
	if err != nil {
		var eo domain.ICalendarRepository
		return eo, err
	}
	o, ok := i.(domain.ICalendarRepository)
	if !ok {
		return o, errors.New("could get 'postgres.repo.calendar' because the object could not be cast to domain.ICalendarRepository")
	}
	return o, nil
}

// GetPostgresRepoCalendar is similar to SafeGetPostgresRepoCalendar but it does not return the error.
// Instead it panics.
func (c *Container) GetPostgresRepoCalendar() domain.ICalendarRepository {
	o, err := c.SafeGetPostgresRepoCalendar()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetPostgresRepoCalendar works like UnscopedSafeGet but only for PostgresRepoCalendar.
// It does not return an interface but a domain.ICalendarRepository.
func (c *Container) UnscopedSafeGetPostgresRepoCalendar() (domain.ICalendarRepository, error) {
	i, err := c.ctn.UnscopedSafeGet("postgres.repo.calendar")
	if err != nil {
		var eo domain.ICalendarRepository
		return eo, err
	}
	o, ok := i.(domain.ICalendarRepository)
	if !ok {
		return o, errors.New("could get 'postgres.repo.calendar' because the object could not be cast to domain.ICalendarRepository")
	}
	return o, nil
}

// UnscopedGetPostgresRepoCalendar is similar to UnscopedSafeGetPostgresRepoCalendar but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetPostgresRepoCalendar() domain.ICalendarRepository {
	o, err := c.UnscopedSafeGetPostgresRepoCalendar()
	if err != nil {
		panic(err)
	}
	return o
}

// PostgresRepoCalendar is similar to GetPostgresRepoCalendar.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetPostgresRepoCalendar method.
// If the container can not be retrieved, it panics.
func PostgresRepoCalendar(i interface{}) domain.ICalendarRepository {
	return C(i).GetPostgresRepoCalendar()
}

// SafeGetPostgresRepoSearch works like SafeGet but only for PostgresRepoSearch.
// It does not return an interface but a domain.ISearchRepository.
func (c *Container) SafeGetPostgresRepoSearch() (domain.ISearchRepository, error) {
	i, err := c.ctn.SafeGet("postgres.repo.search")
	if err != nil {
		var eo domain.ISearchRepository
		return eo, err
	}
	o, ok := i.(domain.ISearchRepository)
	if !ok {
		return o, errors.New("could get 'postgres.repo.search' because the object could not be cast to domain.ISearchRepository")
	}
	return o, nil
}

// GetPostgresRepoSearch is similar to SafeGetPostgresRepoSearch but it does not return the error.
// Instead it panics.
func (c *Container) GetPostgresRepoSearch() domain.ISearchRepository {
	o, err := c.SafeGetPostgresRepoSearch()
	if err != nil {
		panic(err)
	}
	return o
}

// UnscopedSafeGetPostgresRepoSearch works like UnscopedSafeGet but only for PostgresRepoSearch.
// It does not return an interface but a domain.ISearchRepository.
func (c *Container) UnscopedSafeGetPostgresRepoSearch() (domain.ISearchRepository, error) {
	i, err := c.ctn.UnscopedSafeGet("postgres.repo.search")
	if err != nil {
		var eo domain.ISearchRepository
		return eo, err
	}
	o, ok := i.(domain.ISearchRepository)
	if !ok {
		return o, errors.New("could get 'postgres.repo.search' because the object could not be cast to domain.ISearchRepository")
	}
	return o, nil
}

// UnscopedGetPostgresRepoSearch is similar to UnscopedSafeGetPostgresRepoSearch but it does not return the error.
// Instead it panics.
func (c *Container) UnscopedGetPostgresRepoSearch() domain.ISearchRepository {
	o, err := c.UnscopedSafeGetPostgresRepoSearch()
	if err != nil {
		panic(err)
	}
	return o
}

// PostgresRepoSearch is similar to GetPostgresRepoSearch.
// It tries to find the container with the C method and the given interface.
// If the container can be retrieved, it applies the GetPostgresRepoSearch method.
// If the container can not be retrieved, it panics.
func PostgresRepoSearch(i interface{}) domain.ISearchRepository {
	return C(i).GetPostgresRepoSearch()
}
