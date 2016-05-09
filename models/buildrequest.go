package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

/*Buildrequest buildrequest

swagger:model buildrequest
*/
type Buildrequest struct {

	/* git branch checkout
	 */
	Branch string `json:"branch,omitempty"`

	/* build log file
	 */
	Buildlog string `json:"buildlog,omitempty"`

	/* git commit to checkout
	 */
	Commit string `json:"commit,omitempty"`

	/* estimated seconds from now that the build will be completed.  Used to direct when you should check back
	 */
	Eta int32 `json:"eta,omitempty"`

	/* url to the git repo containing a packer config

	Required: true
	*/
	Giturl *string `json:"giturl"`

	/* links to artifacts
	 */
	Images []*Link `json:"images,omitempty"`

	/* not settable
	 */
	Localpath string `json:"localpath,omitempty"`

	/* status of the build
	 */
	Status string `json:"status,omitempty"`

	/* path within the giturl repo to the packer config.  defaults to /packer.json
	 */
	Templatepath string `json:"templatepath,omitempty"`

	/* build log file
	 */
	Testlog string `json:"testlog,omitempty"`
}

// Validate validates this buildrequest
func (m *Buildrequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateGiturl(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateImages(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *Buildrequest) validateGiturl(formats strfmt.Registry) error {

	if err := validate.Required("giturl", "body", m.Giturl); err != nil {
		return err
	}

	return nil
}

func (m *Buildrequest) validateImages(formats strfmt.Registry) error {

	if swag.IsZero(m.Images) { // not required
		return nil
	}

	return nil
}
