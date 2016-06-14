package models

import (
	"fmt"
)

func (m *Buildresponse) String() string {
	return fmt.Sprintf("Buildresponse:[id:%s, requestid:%s, stages:%s]", m.ID, m.Buildrequestid, m.Buildstages)
}
