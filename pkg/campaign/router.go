package campaign

import (
	"image"
	"image/color"
	"image/png"
	"net/http"

	"github.com/tusharsoni/copper/cauth"

	"github.com/gorilla/mux"

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

func NewCreateDraftCampaignRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/campaigns",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.HandleCreateDraftCampaign),
	}}
}

func (ro *Router) HandleCreateDraftCampaign(w http.ResponseWriter, r *http.Request) {
	var body CreateCampaignParams

	if !ro.req.Read(w, r, &body) {
		return
	}

	ctx := r.Context()
	userUUID := cauth.GetCurrentUserUUID(ctx)

	campaign, err := ro.svc.CreateDraftCampaign(ctx, userUUID, body)
	if err != nil {
		ro.logger.Error("Failed to create draft campaign", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, campaign)
}

func NewPublishCampaignRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/campaigns/{uuid}/publish",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.HandlePublishCampaign),
	}}
}

func (ro *Router) HandlePublishCampaign(w http.ResponseWriter, r *http.Request) {
	var campaignUUID = mux.Vars(r)["uuid"]

	err := ro.svc.PublishCampaign(r.Context(), campaignUUID)
	if err != nil {
		ro.logger.Error("Failed to publish campaign", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, nil)
}

func NewTestCampaignRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/campaigns/{uuid}/test",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.HandleTestCampaign),
	}}
}

func (ro *Router) HandleTestCampaign(w http.ResponseWriter, r *http.Request) {
	var (
		body struct {
			Emails []string `json:"emails" valid:"optional"`
		}
		campaignUUID = mux.Vars(r)["uuid"]
	)

	if !ro.req.Read(w, r, &body) {
		return
	}

	err := ro.svc.TestCampaign(r.Context(), campaignUUID, body.Emails)
	if err != nil {
		ro.logger.Error("Failed to send test campaign", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, map[string]bool{
		"success": true,
	})
}

func NewOpenEventImageRoute(ro *Router) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:    "/api/campaigns/{campaignUUID}/events/{contactUUID}/open.png",
		Methods: []string{http.MethodGet},
		Handler: http.HandlerFunc(ro.HandleOpenEventImage),
	}}
}

func (ro *Router) HandleOpenEventImage(w http.ResponseWriter, r *http.Request) {
	const (
		width  = 16
		height = 16
	)

	var (
		campaignUUID = mux.Vars(r)["campaignUUID"]
		contactUUID  = mux.Vars(r)["contactUUID"]
	)

	err := ro.svc.LogEvent(r.Context(), campaignUUID, contactUUID, EventOpen)
	if err != nil {
		ro.logger.WithTags(map[string]interface{}{
			"campaignUUID": campaignUUID,
			"contactUUID":  contactUUID,
			"event":        EventOpen,
		}).Error("Failed to log event", err)
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: 255,
				G: 255,
				B: 255,
				A: 255,
			})
		}
	}

	err = png.Encode(w, img)
	if err != nil {
		ro.logger.Error("Failed to create open event image", err)
		ro.resp.InternalErr(w)
		return
	}
}
