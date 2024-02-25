package muxx

import (
	"net/http"
	"regexp"
)

type Middleware func(http.Handler) http.Handler

type Group struct {
	groupPath   string
	middlewares []Middleware
	mux         *http.ServeMux
}

// New creates a new routes group with a specified base path on top of the existing bundle.
func New() (*Group, error) {
	return &Group{
		mux: http.NewServeMux(),
	}, nil
}

// Mount creates a new routes group with a specified base path on top of the existing bundle.
func Mount(mux *http.ServeMux, basePath string) *Group {
	return &Group{
		mux:       mux,
		groupPath: basePath,
	}
}

// ServeHTTP implements the http.Handler interface
func (b *Group) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	b.mux.ServeHTTP(w, r)
}

// Mount creates a new routes group with a specified base path on top of the existing bundle.
func (b *Group) Mount(groupPath string) *Group {
	middlewares := make([]Middleware, len(b.middlewares))
	copy(middlewares, b.middlewares)
	return &Group{
		mux:         b.mux,
		groupPath:   b.groupPath + groupPath,
		middlewares: middlewares,
	}
}

// Use adds a middleware(s) to the group.
func (b *Group) Use(middlewares ...Middleware) {
	b.middlewares = append(b.middlewares, middlewares...)
}

var httpVerb = regexp.MustCompile(`^(\S*)\s+(.*)$`)

// Handle registers a new handler with the given path and method.
func (b *Group) Handle(path string, handler http.HandlerFunc) {
	w := func(h http.Handler, mws ...Middleware) http.Handler {
		if len(mws) == 0 {
			return h
		}
		res := h
		for i := len(mws) - 1; i >= 0; i-- {
			res = mws[i](res)
		}
		return res
	}

	if b.groupPath != "" {
		matches := httpVerb.FindStringSubmatch(path)
		if len(matches) > 2 {
			path = matches[1] + " " + b.groupPath + matches[2]
		} else {
			path = b.groupPath + path
		}
	}

	b.mux.HandleFunc(path, w(handler, b.middlewares...).ServeHTTP)
}
