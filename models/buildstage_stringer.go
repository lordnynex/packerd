package models

import "fmt"

func (m *Buildstage) String() string {
	return fmt.Sprintf("Buildstage:[ %s, %v, %s]", m.Name, m.Start, m.Status)
	//return "asd"
}
