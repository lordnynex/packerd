package models

// Sanitize should return a build request that is
// safe for display, clearing out passwords and sensitive data
func (br *Buildrequest) Sanitize() *Buildrequest {
	newbr := new(Buildrequest)

	newbr.ID = br.ID
	newbr.Giturl = br.Giturl
	newbr.Branch = br.Branch
	newbr.Responses = br.Responses
	newbr.Templatepath = br.Templatepath
	newbr.Buildonly = br.Buildonly

	return newbr
}

func NewBuildrequest() *Buildrequest {
	br := new(Buildrequest)
	return br
}
