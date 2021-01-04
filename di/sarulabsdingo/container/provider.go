package container

import "github.com/sarulabs/dingo/v4"

// Provider with the test definitions.
type ProviderObject struct {
	dingo.BaseProvider
}

func NewProviderObject() (Provider, error) {
	return &ProviderObject{}, nil
}

func (p *ProviderObject) Load() error {
	for _, fn := range builders {
		if err := fn(p); err != nil {
			return err
		}
	}

	return nil
}
