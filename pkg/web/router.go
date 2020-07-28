package web

import (
	"net/http"

	"github.com/rakyll/statik/fs"
	"github.com/tusharsoni/copper/cerror"
	"github.com/tusharsoni/copper/chttp"

	_ "droplist/pkg/web/build"
)

func NewRouter() (*Router, error) {
	statikFS, err := fs.New()
	if err != nil {
		return nil, cerror.New(err, "failed to create statik fs", nil)
	}

	return &Router{fs: statikFS}, nil
}

type Router struct {
	fs http.FileSystem
}

func NewAppRoute(ro *Router) (chttp.RouteResult, error) {
	index, err := ro.fs.Open("/index.html")
	if err != nil {
		return chttp.RouteResult{}, cerror.New(err, "failed to open index.html", nil)
	}

	indexInfo, err := index.Stat()
	if err != nil {
		return chttp.RouteResult{}, cerror.New(err, "failed to stat index.html", nil)
	}

	return chttp.RouteResult{Route: chttp.Route{
		Path:    "/{path:.*}",
		Methods: []string{http.MethodGet},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeContent(w, r, indexInfo.Name(), indexInfo.ModTime(), index)
		}),
	}}, nil
}

func NewStaticRoute(ro *Router) (chttp.RouteResult, error) {
	return chttp.RouteResult{Route: chttp.Route{
		Path:    "/static/{path:.*}",
		Methods: []string{http.MethodGet},
		Handler: http.FileServer(ro.fs),
	}}, nil
}
