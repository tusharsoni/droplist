package credit

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tusharsoni/copper/cauth"
	"github.com/tusharsoni/copper/chttp"
	"github.com/tusharsoni/copper/clogger"
	"go.uber.org/fx"
)

type Router struct {
	resp   chttp.Responder
	req    chttp.BodyReader
	logger clogger.Logger

	svc    Svc
	config Config
}

type RouterParams struct {
	fx.In

	Resp   chttp.Responder
	Req    chttp.BodyReader
	Logger clogger.Logger

	Svc    Svc
	Config Config
}

func NewRouter(p RouterParams) *Router {
	return &Router{
		resp:   p.Resp,
		req:    p.Req,
		logger: p.Logger,

		svc:    p.Svc,
		config: p.Config,
	}
}

func NewPurchasePackRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/credit/packs/purchase",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.HandlePurchasePack),
	}}
}

func (ro *Router) HandlePurchasePack(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		userUUID = cauth.GetCurrentUserUUID(ctx)
		body     struct {
			ProductID string `json:"product_id" valid:"required"`
		}
	)

	if !ro.req.Read(w, r, &body) {
		return
	}

	pack, secret, err := ro.svc.PurchaseIntent(ctx, userUUID, body.ProductID)
	if err != nil {
		ro.logger.Error("Failed to create pack purchase intent", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, map[string]interface{}{
		"Pack":               pack,
		"StripeClientSecret": secret,
	})
}

func NewCompletePackPurchaseRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/credit/packs/{uuid}/purchase",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.HandleCompletePackPurchase),
	}}
}

func (ro *Router) HandleCompletePackPurchase(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		packUUID = mux.Vars(r)["uuid"]
	)

	err := ro.svc.CompletePurchase(ctx, packUUID)
	if err != nil {
		ro.logger.Error("Failed to complete pack purchase", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, nil)
}

func NewGetValidPacksRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/credit/packs",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodGet},
		Handler:         http.HandlerFunc(ro.HandleGetValidPacks),
	}}
}

func (ro *Router) HandleGetValidPacks(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		userUUID = cauth.GetCurrentUserUUID(ctx)
	)

	packs, err := ro.svc.GetValidPacks(ctx, userUUID)
	if err != nil {
		ro.logger.Error("Failed to get valid packs", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, packs)
}

func NewGetProductsRoute(ro *Router) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:    "/api/credit/products",
		Methods: []string{http.MethodGet},
		Handler: http.HandlerFunc(ro.HandleGetProducts),
	}}
}

func (ro *Router) HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	ro.resp.OK(w, ro.config.Products)
}
