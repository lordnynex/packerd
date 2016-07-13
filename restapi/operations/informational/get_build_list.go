package informational

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetBuildListHandlerFunc turns a function with the right signature into a get build list handler
type GetBuildListHandlerFunc func(GetBuildListParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetBuildListHandlerFunc) Handle(params GetBuildListParams) middleware.Responder {
	return fn(params)
}

// GetBuildListHandler interface for that can handle valid get build list params
type GetBuildListHandler interface {
	Handle(GetBuildListParams) middleware.Responder
}

// NewGetBuildList creates a new http.Handler for the get build list operation
func NewGetBuildList(ctx *middleware.Context, handler GetBuildListHandler) *GetBuildList {
	return &GetBuildList{Context: ctx, Handler: handler}
}

/*GetBuildList swagger:route GET /build/queue informational getBuildList

get a list of links to all build requests

*/
type GetBuildList struct {
	Context *middleware.Context
	Handler GetBuildListHandler
}

func (o *GetBuildList) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewGetBuildListParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}