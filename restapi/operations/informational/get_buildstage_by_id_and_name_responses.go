package informational

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/tompscanlan/packerd/models"
)

/*GetBuildstageByIDAndNameOK list of links to the all stages

swagger:response getBuildstageByIdAndNameOK
*/
type GetBuildstageByIDAndNameOK struct {

	// In: body
	Payload []*models.Link `json:"body,omitempty"`
}

// NewGetBuildstageByIDAndNameOK creates GetBuildstageByIDAndNameOK with default headers values
func NewGetBuildstageByIDAndNameOK() *GetBuildstageByIDAndNameOK {
	return &GetBuildstageByIDAndNameOK{}
}

// WithPayload adds the payload to the get buildstage by Id and name o k response
func (o *GetBuildstageByIDAndNameOK) WithPayload(payload []*models.Link) *GetBuildstageByIDAndNameOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get buildstage by Id and name o k response
func (o *GetBuildstageByIDAndNameOK) SetPayload(payload []*models.Link) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBuildstageByIDAndNameOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if err := producer.Produce(rw, o.Payload); err != nil {
		panic(err) // let the recovery middleware deal with this
	}

}

/*GetBuildstageByIDAndNameBadRequest generic error response

swagger:response getBuildstageByIdAndNameBadRequest
*/
type GetBuildstageByIDAndNameBadRequest struct {

	// In: body
	Payload *models.Error `json:"body,omitempty"`
}

// NewGetBuildstageByIDAndNameBadRequest creates GetBuildstageByIDAndNameBadRequest with default headers values
func NewGetBuildstageByIDAndNameBadRequest() *GetBuildstageByIDAndNameBadRequest {
	return &GetBuildstageByIDAndNameBadRequest{}
}

// WithPayload adds the payload to the get buildstage by Id and name bad request response
func (o *GetBuildstageByIDAndNameBadRequest) WithPayload(payload *models.Error) *GetBuildstageByIDAndNameBadRequest {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get buildstage by Id and name bad request response
func (o *GetBuildstageByIDAndNameBadRequest) SetPayload(payload *models.Error) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetBuildstageByIDAndNameBadRequest) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(400)
	if o.Payload != nil {
		if err := producer.Produce(rw, o.Payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
