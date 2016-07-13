package models

import "fmt"

func (m *Health) String() string {
	return fmt.Sprintf("Health:[ %s, %d]", *m.Status, *m.Diskpercentfull)
}
