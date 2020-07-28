package profile

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/tusharsoni/copper/cerror"

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

func NewGetProfileRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/profile",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodGet},
		Handler:         http.HandlerFunc(ro.HandleGetProfile),
	}}
}

func (ro *Router) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		userUUID = cauth.GetCurrentUserUUID(ctx)
	)

	p, err := ro.svc.GetProfile(ctx, userUUID)
	if err != nil && !cerror.HasCause(err, gorm.ErrRecordNotFound) {
		ro.logger.Error("Failed to get profile", err)
		ro.resp.InternalErr(w)
		return
	}

	if p == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	ro.resp.OK(w, p)
}

func NewSaveProfileRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/profile",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.HandleSaveProfile),
	}}
}

func (ro *Router) HandleSaveProfile(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		userUUID = cauth.GetCurrentUserUUID(ctx)
		body     SaveProfileParams
	)

	if !ro.req.Read(w, r, &body) {
		return
	}

	p, err := ro.svc.SaveProfile(ctx, userUUID, body)
	if err != nil {
		ro.logger.Error("Failed to save profile", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, p)
}
