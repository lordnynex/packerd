package packerd

import (
	"reflect"
	"strings"
	"testing"

	"github.com/tompscanlan/packerd/models"
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
func TestStreamToStringNil(t *testing.T) {
	s := string("")
	StreamToString(strings.NewReader("boing"), &s)
}

func TestDiskCheck(t *testing.T) {
	err, f := GetDirPercentFull("/")
	if err != nil {
		t.Logf("eun expected error found %s", err)
	}
	if f < 0 || f > 100 {
		t.Error("percent out of range")
	}
	t.Logf("got usage: %d", f)
}

func TestDockerShaParse(t *testing.T) {

	var tests = []struct {
		in       string
		expected []string
	}{
		{`Imported Docker image: sha256:e8d78bc014ef77b9031de45fc9dfb7f5e721d14a3b7cdd0816c5d3751c7b227b`,
			[]string{"e8d78bc014ef77b9031de45fc9dfb7f5e721d14a3b7cdd0816c5d3751c7b227b"}},
		{
			"no sha here",
			[]string{},
		},
		{"d009e209a0c4",
			[]string{}},
		{`2016/06/20 15:23:32 packer: 2016/06/20 15:23:32 Executing in container ca55429e8abac80ecbdf5b00368943372a9e587da77d7205c8df9745406b9027: "(rm -f /tmp/script_4051.sh) >/packer-files/cmd078235457 2>&1; echo $? >/packer-files/cmd078235457-exit"`,
			[]string{}},
		{`2016/06/20 15:29:14 ui: --> ubuntu-docker: Imported Docker image: sha256:e8d78bc014ef77b9031de45fc9dfb7f5e721d14a3b7cdd0816c5d3751c7b227b`,
			[]string{"e8d78bc014ef77b9031de45fc9dfb7f5e721d14a3b7cdd0816c5d3751c7b227b"}},

		{`2016/06/20 15:29:14 machine readable: ubuntu-docker,artifact []string{"0", "builder-id", "packer.post-processor.docker-import"}
2016/06/20 15:29:14 machine readable: ubuntu-docker,artifact []string{"0", "id", "sha256:e8d78bc014ef77b9031de45fc9dfb7f5e721d14a3b7cdd0816c5d3751c7b227b"}
2016/06/20 15:29:14 machine readable: ubuntu-docker,artifact []string{"0", "string", "Imported Docker image: sha256:e8d78bc014ef77b9031de45fc9dfb7f5e721d14a3b7cdd0816c5d3751c7b227b"}`,
			[]string{"e8d78bc014ef77b9031de45fc9dfb7f5e721d14a3b7cdd0816c5d3751c7b227b"}},
	}

	for _, test := range tests {

		actual := ParseDockerImageSha(test.in)

		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("ParseDockerImageSha(%s): expected [%s], got [%s]", test.in, test.expected, actual)
		}
	}
}
func TestDockerUrlParse(t *testing.T) {

	var tests = []struct {
		in       string
		expected []string
	}{
		{
			"no url here",
			[]string{},
		},
		{"d009e209a0c4",
			[]string{}},
		{`2016/06/20 15:29:05 ui: ubuntu-docker (docker-tag): Tagging image: sha256:e8d78bc014ef77b9031de45fc9dfb7f5e721d14a3b7cdd0816c5d3751c7b227b
2016/06/20 15:29:05 ui: ubuntu-docker (docker-tag): Repository: dockyardaws.cloud.capitalone.com/tompscanlan/packerd:latest`,
			[]string{"dockyardaws.cloud.capitalone.com/tompscanlan/packerd:latest"}},
		{`2016/06/20 15:24:59 ui: ubuntu-docker (docker-tag): Tagging image: sha256:e8d78bc014ef77b9031de45fc9dfb7f5e721d14a3b7cdd0816c5d3751c7b227b
2016/06/20 15:24:59 ui: ubuntu-docker (docker-tag): Repository: dockyardaws.cloud.capitalone.com/tompscanlan/packerd:0.135
2016/06/20 15:24:59 Flagging to keep original artifact from post-processor 'docker-tag'`,
			[]string{"dockyardaws.cloud.capitalone.com/tompscanlan/packerd:0.135"},
		}}

	for _, test := range tests {
		actual := ParseDockerArtifactUrl(test.in)

		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("ParseDockerArtifactUrl(%s): expected [%s], got [%s]", test.in, test.expected, actual)
		}
	}
}
