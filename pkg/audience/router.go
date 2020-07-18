package audience

import (
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

func NewCreateListRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/audience/lists",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.HandleCreateList),
	}}
}

func (ro *Router) HandleCreateList(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name" valid:"required"`
	}

	if !ro.req.Read(w, r, &body) {
		return
	}

	ctx := r.Context()
	userUUID := cauth.GetCurrentUserUUID(ctx)

	list, err := ro.svc.CreateList(ctx, body.Name, userUUID)
	if err != nil {
		ro.logger.Error("Failed to create list", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, list)
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

	ctx := r.Context()
	userUUID := cauth.GetCurrentUserUUID(ctx)

	contacts, err := ro.svc.ListUserContacts(ctx, userUUID)
	if err != nil {
		ro.logger.Error("Failed to list user contacts", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, contacts)
}

func NewAddContactsToListRoute(ro *Router, auth cauth.Middleware) chttp.RouteResult {
	return chttp.RouteResult{Route: chttp.Route{
		Path:            "/api/audience/lists/{uuid}/contacts",
		MiddlewareFuncs: []chttp.MiddlewareFunc{auth.VerifySessionToken},
		Methods:         []string{http.MethodPost},
		Handler:         http.HandlerFunc(ro.HandleAddContactsToList),
	}}
}

func (ro *Router) HandleAddContactsToList(w http.ResponseWriter, r *http.Request) {
	var (
		body struct {
			ContactUUIDs []string `json:"contact_uuids" valid:"optional"`
		}
		listUUID = mux.Vars(r)["uuid"]
	)

	if !ro.req.Read(w, r, &body) {
		return
	}

	results, err := ro.svc.AddContactsToList(r.Context(), listUUID, body.ContactUUIDs)
	if err != nil {
		ro.logger.Error("Failed to add contacts to list", err)
		ro.resp.InternalErr(w)
		return
	}

	ro.resp.OK(w, results)
}
