package packerd

type Worker interface {
	Start()
	Stop()
}

// potential states of a worker
const (
	Stopped  = 0
	Stopping = 1
	Started  = 2
)
