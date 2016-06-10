package packerd

import (
	"testing"

	"github.kdc.capitalone.com/kbs316/packerd/models"
)

func strPtr(s string) *string { return &s }

var flagtests = []struct {
	in      models.Buildrequest
	out     string
	invalid bool
}{
	{
		models.Buildrequest{
			Giturl: nil,
		},
		`{"giturl":null}`,
		true,
	},
	{
		models.Buildrequest{
			Giturl: strPtr("http://github.com/"),
		},
		`{"giturl":"http://github.com/"}`,
		false,
	},
	{
		models.Buildrequest{
			Giturl: strPtr("foo"),
		},
		`{"giturl":"foo"}`,
		false,
	},
}

func TestPrintBuildRequest(t *testing.T) {

	for _, tt := range flagtests {

		found := BuildRequestToString(tt.in)

		err := tt.in.Validate(nil)
		if err == nil && tt.invalid == true {
			t.Errorf("expected invalid to be %t, found %q", tt.invalid, err)
		}

		if found != tt.out {
			t.Errorf("expected %q, found %q", tt.out, found)
		}
	}
}
