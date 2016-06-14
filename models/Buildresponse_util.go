package models

func NewBuildresponse() *Buildresponse {
	resp := new(Buildresponse)
	resp.Status = PENDING

	return resp
}
