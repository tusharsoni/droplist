package audience

import (
	"errors"
	"net/http"
	"strconv"

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

func NewSummaryRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/audience/summary",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodGet},
		Handler:         http.HandlerFunc(ro.HandleSummary),
	}}
}

func (ro *Router) HandleSummary(w http.ResponseWriter, r *http.Request) {
	var (
		ctx      = r.Context()
		userUUID = cauth.GetCurrentUserUUID(ctx)
	)

	summary, err := ro.svc.Summary(ctx, userUUID)
	if err != nil {
		ro.logger.Error("Failed to get audience summary", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, summary)
}

func NewCreateContactsRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/audience/contacts",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.HandleCreateContacts),
	}}
}

func (ro *Router) HandleCreateContacts(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Contacts []CreateContactParams `json:"contacts" valid:"optional"`
	}

	if !ro.req.Read(w, r, &body) {
		return
	}

	ctx := r.Context()
	userUUID := cauth.GetCurrentUserUUID(ctx)

	results, err := ro.svc.CreateContacts(ctx, userUUID, body.Contacts)
	if err != nil {
		ro.logger.Error("Failed to create contacts", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, results)
}

func NewListContactsRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/audience/contacts",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodGet},
		Handler:         http.HandlerFunc(ro.HandleListContacts),
	}}
}

func (ro *Router) HandleListContacts(w http.ResponseWriter, r *http.Request) {
	const (
		defaultLimit  = int(20)
		defaultOffset = int(0)
	)

	var (
		ctx      = r.Context()
		userUUID = cauth.GetCurrentUserUUID(ctx)

		rawLimit  = r.URL.Query().Get("limit")
		rawOffset = r.URL.Query().Get("offset")
	)

	limit, err := strconv.Atoi(rawLimit)
	if err != nil {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(rawOffset)
	if err != nil {
		offset = defaultOffset
	}

	contacts, err := ro.svc.ListUserContacts(ctx, userUUID, limit, offset)
	if err != nil {
		ro.logger.Error("Failed to list user contacts", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, contacts)
}

func NewDeleteContactsRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/audience/contacts",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodDelete},
		Handler:         http.HandlerFunc(ro.HandleDeleteContacts),
	}}
}

func (ro *Router) HandleDeleteContacts(w http.ResponseWriter, r *http.Request) {
	var (
		body struct {
			DeleteAll    bool     `json:"delete_all" valid:"optional"`
			ContactUUIDs []string `json:"contact_uuids" valid:"optional"`
		}
		ctx      = r.Context()
		userUUID = cauth.GetCurrentUserUUID(ctx)
	)

	if !ro.req.Read(w, r, &body) {
		return
	}

	if !body.DeleteAll && len(body.ContactUUIDs) == 0 {
		ro.resp.BadRequest(w, errors.New("contact uuids are required if delete all is false"))
		return
	}

	if body.DeleteAll {
		body.ContactUUIDs = nil
	}

	err := ro.svc.DeleteContacts(ctx, userUUID, body.ContactUUIDs)
	if err != nil {
		ro.logger.Error("Failed to delete contacts", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, nil)
}

func NewUnsubscribeContactRoute(ro *Router) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:    "/api/audience/contacts/{uuid}/unsubscribe",
		Methods: []string{http.MethodGet},
		Handler: http.HandlerFunc(ro.HandleUnsubscribeContact),
	}}
}

func (ro *Router) HandleUnsubscribeContact(w http.ResponseWriter, r *http.Request) {
	var contactUUID = mux.Vars(r)["uuid"]

	err := ro.svc.UnsubscribeContact(r.Context(), contactUUID)
	if err != nil {
		ro.logger.Error("Failed to unsubscribe contact", err)
		ro.resp.InternalErr(w)
		return
	}

	_, err = w.Write([]byte("You have been unsubscribed from this mailing list."))
	if err != nil {
		ro.logger.Error("Failed to write response to body", err)
	}
}
