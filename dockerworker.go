package packerd

import (
	"strings"

	log "github.com/Sirupsen/logrus"

	"github.com/tompscanlan/packerd/models"
)

type DockerWorker struct {
	ID    int
	Done  chan bool
	State int

	// as above, for docker clean work
	Work      chan *DockerClean
	WorkQueue chan chan *DockerClean
}

// NewDockerWorker makes a new worker
func NewDockerWorker(id int, q chan chan *DockerClean) *DockerWorker {

	worker := new(DockerWorker)

	worker.ID = id
	worker.Done = make(chan bool)
	worker.Work = make(chan *DockerClean)
	worker.WorkQueue = q

	worker.State = Stopped

	return worker
}

//  RunDockerCleanup removes unnedded containers and images
func (w *DockerWorker) RunDockerCleanup(dc *DockerClean) error {
	stage := models.NewBuildstage("DockerCleanup")
	args := []string{"images", "-f", "dangling=true", "-q"}
	err := RunCmd("docker", args, "/tmp", nil, stage)

	args = []string{"rmi"}
	if stage.Log != "" {
		args = append(args, strings.Split(stage.Log, "\n")...)
		err := RunCmd("docker", args, "/tmp", nil, stage)

		if err != nil {
			if strings.Contains(err.Error(), "must be forced") {
				log.Error("Some docker images could not be cleaned. Clean containers referencing them to resolve.")
			} else {
				return err
			}
		}
	}

	return err
}

// Start gets the worker recieving from the global WorkQueue and running the needed stages
func (w *DockerWorker) Start() {
	w.State = Started

	go func() {
		for {
			// put our work channel onto the dispatcher's queue.
			// dispatcher will shift it off, and send work down it
			// Re-add it after doing the job
			w.WorkQueue <- w.Work

			select {
			// handle a docker clean request
			case clean := <-w.Work:
				if ProccessNameExists("kitchen") {
					log.Info("Skipping docker clean while kitchen is running")
				} else if ProccessNameExists("packer") {
					log.Info("Skipping docker clean while packer is running")
				} else {
					log.Infof("Starting a docker clean job %v", clean)
					if err := w.RunDockerCleanup(clean); err != nil {
						log.Error(err.Error())
					}
					log.Debug("after docker clean")
				}

			case <-w.Done:
				log.Debugf("worker%d: done", w.ID)
				return
			}
		}
	}()
}

// Stop this worker
func (w *DockerWorker) Stop() {
	w.State = Stopping
	go func() {
		w.Done <- true
		w.State = Stopped
	}()
}
