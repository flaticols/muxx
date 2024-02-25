package muxx

import (
	"errors"
	"net/http"
)

type Option func(*Options) error

type Options struct {
	Mux *http.ServeMux
}

func WithMux(mux *http.ServeMux) Option {
	return func(o *Options) error {
		o.Mux = mux
		return nil
	}
}

func (o *Options) Apply(opts ...Option) error {
	var err error
	for _, opt := range opts {
		err := opt(o)
		if err != nil {
			err = errors.Join(err, err)
		}
	}

	return err
}

func New(rootPath string, options ...Option) (*Group, error) {
	o := &Options{
		Mux: http.NewServeMux(),
	}
	err := o.Apply(options...)
	if err != nil {
		return nil, err
	}

	return &Group{
		mux:       o.Mux,
		groupPath: rootPath,
	}, nil
}
