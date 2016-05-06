package informational

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/validate"

	strfmt "github.com/go-openapi/strfmt"
)

// NewGetQueueTestLogByIDParams creates a new GetQueueTestLogByIDParams object
// with the default values initialized.
func NewGetQueueTestLogByIDParams() GetQueueTestLogByIDParams {
	var ()
	return GetQueueTestLogByIDParams{}
}

// GetQueueTestLogByIDParams contains all the bound params for the get queue test log by Id operation
// typically these are obtained from a http.Request
//
// swagger:parameters getQueueTestLogById
type GetQueueTestLogByIDParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request

	/*uuid for the build
	  Required: true
	  Max Length: 36
	  Min Length: 36
	  In: path
	*/
	ID string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls
func (o *GetQueueTestLogByIDParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error
	o.HTTPRequest = r

	rID, rhkID, _ := route.Params.GetOK("id")
	if err := o.bindID(rID, rhkID, route.Formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (o *GetQueueTestLogByIDParams) bindID(rawData []string, hasKey bool, formats strfmt.Registry) error {
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	o.ID = raw

	if err := o.validateID(formats); err != nil {
		return err
	}

	return nil
}

func (o *GetQueueTestLogByIDParams) validateID(formats strfmt.Registry) error {

	if err := validate.MinLength("id", "path", string(o.ID), 36); err != nil {
		return err
	}

	if err := validate.MaxLength("id", "path", string(o.ID), 36); err != nil {
		return err
	}

	return nil
}
