// dispatcher
package packerd

import (
	log "github.com/Sirupsen/logrus"

	"github.com/tompscanlan/packerd/models"
)

// for sending build requests to the worker pool
var WorkQueue = make(chan *models.Buildrequest, 100)

// channel to pass work channel to worker
var WorkerQueue chan chan *models.Buildrequest

func StartDispatcher(count int) {

	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan *models.Buildrequest, count)

	for i := 1; i < count+1; i++ {
		log.Debugf("Dispatcher starting worker%d", i)
		worker := NewWorker(i, WorkerQueue)
		worker.Start()
	}

	go func() {
		for {
			select {
			case work := <-WorkQueue:
				log.Debugf("Received build requeust", BuildRequestToString(*work))
				go func() {
					worker := <-WorkerQueue

					log.Debugln("Dispatching build request", BuildRequestToString(*work))
					worker <- work

				}()
			}
		}
	}()
}
