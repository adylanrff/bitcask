package bitcask

import (
	"errors"

	"github.com/adylanrff/bitcask/internal/minimal"
	"github.com/adylanrff/bitcask/pkg/option"
	"github.com/adylanrff/bitcask/pkg/types"
)

func Open(dirname string, opts ...option.OptionFunc) (types.BitcaskHandler, error) {
	if err := validateOpenArguments(dirname, opts...); err != nil {
		return nil, err
	}

	options := option.NewDefaultOptions()
	options.Apply(opts...)

	handler, err := minimal.NewHandler(dirname, options)
	if err != nil {
		return nil, err
	}

	if err := handler.Init(); err != nil {
		return nil, err
	}

	return handler, nil
}

func validateOpenArguments(dirname string, opts ...option.OptionFunc) error {
	if dirname == "" {
		return errors.New("dirname cannot be empty")
	}
	return nil
}
