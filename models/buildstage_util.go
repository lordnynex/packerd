package models

// states of a build stage
const (
	PENDING  = "pending"
	RUNNING  = "running"
	FAILED   = "failed"
	COMPLETE = "complete"
)

func NewBuildstage(name string) *Buildstage {
	stage := new(Buildstage)
	stage.Name = name
	stage.Status = PENDING

	return stage
}
