package models

import (
	"fmt"
)

func (m *Buildrequest) String() string {
	return fmt.Sprintf("Buildrequest:[id:%s, giturl:%s, branch:%s, templatepath:%s, responses:%s]", m.ID, *m.Giturl, m.Branch, m.Templatepath, m.Responses)
}
