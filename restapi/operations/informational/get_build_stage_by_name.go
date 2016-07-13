package informational

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetBuildStageByNameHandlerFunc turns a function with the right signature into a get build stage by name handler
type GetBuildStageByNameHandlerFunc func(GetBuildStageByNameParams) middleware.Responder

// Handle executing the request and returning a response
func (fn GetBuildStageByNameHandlerFunc) Handle(params GetBuildStageByNameParams) middleware.Responder {
	return fn(params)
}

// GetBuildStageByNameHandler interface for that can handle valid get build stage by name params
type GetBuildStageByNameHandler interface {
	Handle(GetBuildStageByNameParams) middleware.Responder
}

// NewGetBuildStageByName creates a new http.Handler for the get build stage by name operation
func NewGetBuildStageByName(ctx *middleware.Context, handler GetBuildStageByNameHandler) *GetBuildStageByName {
	return &GetBuildStageByName{Context: ctx, Handler: handler}
}

/*GetBuildStageByName swagger:route GET /build/responses/{id}/{buildnumber}/stages/{stagename} informational getBuildStageByName

get all the stages of a specific build response

*/
type GetBuildStageByName struct {
	Context *middleware.Context
	Handler GetBuildStageByNameHandler
}

func (o *GetBuildStageByName) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)
	var Params = NewGetBuildStageByNameParams()

	if err := o.Context.BindValidRequest(r, route, &Params); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle(Params) // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}