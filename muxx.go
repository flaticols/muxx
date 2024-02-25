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
func (g *Group) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	g.mux.ServeHTTP(w, r)
}

// Mount creates a new routes group with a specified base path on top of the existing bundle.
func (g *Group) Mount(groupPath string) *Group {
	middlewares := make([]Middleware, len(g.middlewares))
	copy(middlewares, g.middlewares)
	return &Group{
		mux:         g.mux,
		groupPath:   g.groupPath + groupPath,
		middlewares: middlewares,
	}
}

// Use adds a middleware(s) to the group.
func (g *Group) Use(middlewares ...Middleware) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *Group) Group() *Group {
	middlewares := make([]Middleware, len(g.middlewares))
	copy(middlewares, g.middlewares)
	return &Group{
		mux:         g.mux,
		groupPath:   g.groupPath,
		middlewares: middlewares,
	}
}

func (g *Group) Route(fn func(*Group)) {
	fn(g)
}

var httpVerb = regexp.MustCompile(`^(\S*)\s+(.*)$`)

// Handle registers a new handler with the given path and method.
func (g *Group) Handle(path string, handler http.HandlerFunc) {
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

	if g.groupPath != "" {
		matches := httpVerb.FindStringSubmatch(path)
		if len(matches) > 2 {
			path = matches[1] + " " + g.groupPath + matches[2]
		} else {
			path = g.groupPath + path
		}
	}

	g.mux.HandleFunc(path, w(handler, g.middlewares...).ServeHTTP)
}
