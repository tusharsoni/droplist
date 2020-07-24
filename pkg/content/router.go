package content

import (
	"net/http"

	"github.com/tusharsoni/copper/cauth"

	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
)

type Router struct {
	resp   chttp.Responder
	req    chttp.BodyReader
	logger clogger.Logger

	svc Svc
}

type RouterParams struct {
	fx.In

	Resp   chttp.Responder
	Req    chttp.BodyReader
	Logger clogger.Logger

	Svc Svc
}

func NewRouter(p RouterParams) *Router {
	return &Router{
		resp:   p.Resp,
		req:    p.Req,
		logger: p.Logger,

		svc: p.Svc,
	}
}

func NewCreateTemplateRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/content/templates",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.HandleCreateTemplate),
	}}
}

func (ro *Router) HandleCreateTemplate(w http.ResponseWriter, r *http.Request) {
	var body CreateTemplateParams

	if !ro.req.Read(w, r, &body) {
		return
	}

	ctx := r.Context()
	userUUID := cauth.GetCurrentUserUUID(ctx)

	tmpl, err := ro.svc.CreateTemplate(ctx, userUUID, body)
	if err != nil {
		ro.logger.Error("Failed to create template", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, tmpl)
}

func NewListTemplatesRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/content/templates",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodGet},
		Handler:         http.HandlerFunc(ro.HandleListTemplates),
	}}
}

func (ro *Router) HandleListTemplates(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		userUUID = cauth.GetCurrentUserUUID(ctx)
	)

	templates, err := ro.svc.ListUserTemplates(ctx, userUUID)
	if err != nil {
		ro.logger.Error("Failed to list user templates", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, templates)
}
