package container

import (
	"fmt"
	"sync"

	"github.com/sarulabs/di/v2"
)

var (
	mu       sync.Mutex
	context  di.Container
	builders []buildFn
)

type (
	// public
	Tag        = di.Tag
	Builder    = di.Builder
	Context    = di.Container
	Definition = di.Def

	// private
	buildFn func(builder *di.Builder) error
)

// Register definition builder
func Register(fn buildFn) {
	mu.Lock()
	defer mu.Unlock()

	builders = append(builders, fn)
}

// Get context
func Instance() (di.Container, error) {
	if context != nil {
		return context, nil
	}

	builder, err := di.NewBuilder()
	if err != nil {
		return nil, fmt.Errorf("can't create context builder: %s", err)
	}

	for _, fn := range builders {
		if err := fn(builder); err != nil {
			return nil, err
		}
	}

	context = builder.Build()

	return context, nil
}
