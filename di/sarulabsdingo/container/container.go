package container

import (
	"sync"

	"github.com/sarulabs/di/v2"
	"github.com/sarulabs/dingo/v4"
)

var (
	mu       sync.Mutex
	context  di.Container
	builders []buildFn
)

type (
	// public
	Params     = dingo.Params
	Service    = dingo.Service
	Definition = dingo.Def
	Provider   = dingo.Provider
	Container  = di.Container

	// private
	buildFn func(builder *ProviderObject) error
)

// Register definition builder
func Register(fn buildFn) {
	mu.Lock()
	defer mu.Unlock()

	builders = append(builders, fn)
}
