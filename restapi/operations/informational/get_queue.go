package informational

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the generate command

import (
	"net/http"

	middleware "github.com/go-openapi/runtime/middleware"
)

// GetQueueHandlerFunc turns a function with the right signature into a get queue handler
type GetQueueHandlerFunc func() middleware.Responder

// Handle executing the request and returning a response
func (fn GetQueueHandlerFunc) Handle() middleware.Responder {
	return fn()
}

// GetQueueHandler interface for that can handle valid get queue params
type GetQueueHandler interface {
	Handle() middleware.Responder
}

// NewGetQueue creates a new http.Handler for the get queue operation
func NewGetQueue(ctx *middleware.Context, handler GetQueueHandler) *GetQueue {
	return &GetQueue{Context: ctx, Handler: handler}
}

/*GetQueue swagger:route GET /queue informational getQueue

get a list of links to all build status

*/
type GetQueue struct {
	Context *middleware.Context
	Handler GetQueueHandler
}

func (o *GetQueue) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	route, _ := o.Context.RouteInfo(r)

	if err := o.Context.BindValidRequest(r, route, nil); err != nil { // bind params
		o.Context.Respond(rw, r, route.Produces, route, err)
		return
	}

	res := o.Handler.Handle() // actually handle the request

	o.Context.Respond(rw, r, route.Produces, route, res)

}
