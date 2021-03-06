package informational

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/tompscanlan/packerd/models"
)

/*GetBuildResponsesIDBuildnumberOK returns a build status

swagger:response getBuildResponsesIdBuildnumberOK
*/
type GetBuildResponsesIDBuildnumberOK struct {

	// In: body
	Payload *models.Buildresponse `json:"body,omitempty"`
}

// NewGetBuildResponsesIDBuildnumberOK creates GetBuildResponsesIDBuildnumberOK with default headers values
func NewGetBuildResponsesIDBuildnumberOK() *GetBuildResponsesIDBuildnumberOK {
	return &GetBuildResponsesIDBuildnumberOK{}
}

// WithPayload adds the payload to the get build responses Id buildnumber o k response
func (o *GetBuildResponsesIDBuildnumberOK) WithPayload(payload *models.Buildresponse) *GetBuildResponsesIDBuildnumberOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get build responses Id buildnumber o k response
func (o *GetBuildResponsesIDBuildnumberOK) SetPayload(payload *models.Buildresponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBuildResponsesIDBuildnumberOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

/*GetBuildResponsesIDBuildnumberBadRequest generic error response

swagger:response getBuildResponsesIdBuildnumberBadRequest
*/
type GetBuildResponsesIDBuildnumberBadRequest struct {

	// In: body
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetBuildResponsesIDBuildnumberBadRequest creates GetBuildResponsesIDBuildnumberBadRequest with default headers values
func NewGetBuildResponsesIDBuildnumberBadRequest() *GetBuildResponsesIDBuildnumberBadRequest {
	return &GetBuildResponsesIDBuildnumberBadRequest{}
}

// WithPayload adds the payload to the get build responses Id buildnumber bad request response
func (o *GetBuildResponsesIDBuildnumberBadRequest) WithPayload(payload *models.Error) *GetBuildResponsesIDBuildnumberBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get build responses Id buildnumber bad request response
func (o *GetBuildResponsesIDBuildnumberBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBuildResponsesIDBuildnumberBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
