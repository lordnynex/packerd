package informational

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetBuildStagesByIDBuildNumberHandlerFunc turns a function with the right signature into a get build stages by Id build number handler
type GetBuildStagesByIDBuildNumberHandlerFunc func(GetBuildStagesByIDBuildNumberParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetBuildStagesByIDBuildNumberHandlerFunc) Handle(params GetBuildStagesByIDBuildNumberParams) middleware.Responder {
	return fn(params)
}

// GetBuildStagesByIDBuildNumberHandler interface for that can handle valid get build stages by Id build number params
type GetBuildStagesByIDBuildNumberHandler interface {
	Handle(GetBuildStagesByIDBuildNumberParams) middleware.Responder
}

// NewGetBuildStagesByIDBuildNumber creates a new http.Handler for the get build stages by Id build number operation
func NewGetBuildStagesByIDBuildNumber(ctx *middleware.Context, handler GetBuildStagesByIDBuildNumberHandler) *GetBuildStagesByIDBuildNumber {
	return &GetBuildStagesByIDBuildNumber{Context: ctx, Handler: handler}
}

/*GetBuildStagesByIDBuildNumber swagger:route GET /build/responses/{id}/{buildnumber}/stages informational getBuildStagesByIdBuildNumber

get all the stages of a specific build response

*/
type GetBuildStagesByIDBuildNumber struct {
	Context *middleware.Context
	Handler GetBuildStagesByIDBuildNumberHandler
}

func (o *GetBuildStagesByIDBuildNumber) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewGetBuildStagesByIDBuildNumberParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}