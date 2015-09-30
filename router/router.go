package router

import (
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

var (
	mutex = sync.RWMutex{}
	ctxts = map[*http.Request]context.Context{}
)

// Context will return the context for the request.
func Context(r *http.Request) context.Context {
	mutex.RLock()
	defer mutex.RUnlock()
	return ctxts[r]
}

type param string

// Param returns param p for the context.
func Param(ctx context.Context, p string) string {
	return ctx.Value(param(p)).(string)
}

// handle turns a Handle into httprouter.Handle
func handle(handlerFunc http.HandlerFunc) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		for _, p := range params {
			ctx = context.WithValue(ctx, param(p.Key), p.Value)
		}

		mutex.Lock()
		ctxts[r] = ctx
		mutex.Unlock()

		handlerFunc(w, r)

		mutex.Lock()
		delete(ctxts, r)
		mutex.Unlock()
	}
}

// Router is wrapper around httprouter.Router with added support for prefixed sub-routes.
type Router struct {
	rtr    *httprouter.Router
	prefix string
}

// New returns a new Router.
func New() *Router {
	return &Router{rtr: httprouter.New()}
}

// WithPrefix returns a router that prefixes all registered routes with prefix.
func (r *Router) WithPrefix(prefix string) *Router {
	return &Router{rtr: r.rtr, prefix: r.prefix + prefix}
}

// Map all supported httprouter methods

// Del registers a new DELETE route.
func (r *Router) Del(path string, h http.HandlerFunc) {
	r.rtr.DELETE(r.prefix+path, handle(h))
}

// Get registers a new GET route.
func (r *Router) Get(path string, h http.HandlerFunc) {
	r.rtr.GET(r.prefix+path, handle(h))
}

// Head registers a new HEAD route.
func (r *Router) Head(path string, h http.HandlerFunc) {
	r.rtr.HEAD(r.prefix+path, handle(h))
}

// Options registers a new OPTIONS route.
func (r *Router) Options(path string, h http.HandlerFunc) {
	r.rtr.OPTIONS(r.prefix+path, handle(h))
}

// Patch registers a new PATCH route.
func (r *Router) Patch(path string, h http.HandlerFunc) {
	r.rtr.PATCH(r.prefix+path, handle(h))
}

// Post registers a new POST route.
func (r *Router) Post(path string, h http.HandlerFunc) {
	r.rtr.POST(r.prefix+path, handle(h))
}

// Put registers a new PUT route.
func (r *Router) Put(path string, h http.HandlerFunc) {
	r.rtr.PUT(r.prefix+path, handle(h))
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.rtr.ServeHTTP(w, req)
}

// FileServe returns a new http.HandlerFunc that serves files from dir.
// Using routes must provide the *filepath parameter.
func (r *Router) FileServe(dir string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))

	return func(w http.ResponseWriter, req *http.Request) {
		req.URL.Path = Param(Context(req), "filepath")
		fs.ServeHTTP(w, req)
	}
}
