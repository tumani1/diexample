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
	buildFn func(builder *di.Builder, params map[string]interface{}) error
)

// Register definition builder
func Register(fn buildFn) {
	mu.Lock()
	defer mu.Unlock()

	builders = append(builders, fn)
}

// Get context
func Instance(params map[string]interface{}) (di.Container, error) {
	if context != nil {
		return context, nil
	}

	builder, err := di.NewBuilder()
	if err != nil {
		return nil, fmt.Errorf("can't create context builder: %s", err)
	}

	for _, fn := range builders {
		if err := fn(builder, params); err != nil {
			return nil, err
		}
	}

	context = builder.Build()

	return context, nil
}

// Iterate definitions by tag
func Iterate(ctx di.Container, tag string, fn func(ctx Context, tag *Tag, name string) error) (err error) {
	for name, def := range ctx.Definitions() {
		for _, defTag := range def.Tags {
			if defTag.Name != tag {
				continue
			}

			if err = fn(ctx, &defTag, name); err != nil {
				return err
			}

			break
		}
	}

	return nil
}
