package informational

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/tompscanlan/packerd/models"
)

/*GetBuildstagesByIDOK list of links to the all stages

swagger:response getBuildstagesByIdOK
*/
type GetBuildstagesByIDOK struct {

	// In: body
	Payload []*models.Buildstage `json:"body,omitempty"`
}

// NewGetBuildstagesByIDOK creates GetBuildstagesByIDOK with default headers values
func NewGetBuildstagesByIDOK() *GetBuildstagesByIDOK {
	return &GetBuildstagesByIDOK{}
}

// WithPayload adds the payload to the get buildstages by Id o k response
func (o *GetBuildstagesByIDOK) WithPayload(payload []*models.Buildstage) *GetBuildstagesByIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get buildstages by Id o k response
func (o *GetBuildstagesByIDOK) SetPayload(payload []*models.Buildstage) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBuildstagesByIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

/*GetBuildstagesByIDBadRequest generic error response

swagger:response getBuildstagesByIdBadRequest
*/
type GetBuildstagesByIDBadRequest struct {

	// In: body
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetBuildstagesByIDBadRequest creates GetBuildstagesByIDBadRequest with default headers values
func NewGetBuildstagesByIDBadRequest() *GetBuildstagesByIDBadRequest {
	return &GetBuildstagesByIDBadRequest{}
}

// WithPayload adds the payload to the get buildstages by Id bad request response
func (o *GetBuildstagesByIDBadRequest) WithPayload(payload *models.Error) *GetBuildstagesByIDBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get buildstages by Id bad request response
func (o *GetBuildstagesByIDBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBuildstagesByIDBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
