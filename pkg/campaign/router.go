package campaign

import (
	"net/http"

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

func NewCreateDraftCampaignRoute(ro *Router) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:    "/api/campaigns/create",
		Methods: []string{http.MethodPost},
		Handler: http.HandlerFunc(ro.HandleCreateDraftCampaign),
	}}
}

func (ro *Router) HandleCreateDraftCampaign(w http.ResponseWriter, r *http.Request) {
	var body CreateCampaignParams

	if !ro.req.Read(w, r, &body) {
		return
	}

	campaign, err := ro.svc.CreateDraftCampaign(r.Context(), body)
	if err != nil {
		ro.logger.Error("Failed to create draft campaign", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, campaign)
}
