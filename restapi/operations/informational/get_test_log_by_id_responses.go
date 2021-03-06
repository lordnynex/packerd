package informational

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/tompscanlan/packerd/models"
)

/*GetTestLogByIDOK returns the test kitchen run log

swagger:response getTestLogByIdOK
*/
type GetTestLogByIDOK struct {

	// In: body
	Payload string `json:"body,omitempty"`
}

// NewGetTestLogByIDOK creates GetTestLogByIDOK with default headers values
func NewGetTestLogByIDOK() *GetTestLogByIDOK {
	return &GetTestLogByIDOK{}
}

// WithPayload adds the payload to the get test log by Id o k response
func (o *GetTestLogByIDOK) WithPayload(payload string) *GetTestLogByIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get test log by Id o k response
func (o *GetTestLogByIDOK) SetPayload(payload string) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTestLogByIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

/*GetTestLogByIDBadRequest generic error response

swagger:response getTestLogByIdBadRequest
*/
type GetTestLogByIDBadRequest struct {

	// In: body
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetTestLogByIDBadRequest creates GetTestLogByIDBadRequest with default headers values
func NewGetTestLogByIDBadRequest() *GetTestLogByIDBadRequest {
	return &GetTestLogByIDBadRequest{}
}

// WithPayload adds the payload to the get test log by Id bad request response
func (o *GetTestLogByIDBadRequest) WithPayload(payload *models.Error) *GetTestLogByIDBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get test log by Id bad request response
func (o *GetTestLogByIDBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetTestLogByIDBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
