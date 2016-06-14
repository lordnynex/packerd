package packerd

import (
	log "github.com/Sirupsen/logrus"

	"github.com/tompscanlan/packerd/models"
)

// BuildSendWorkerChan is for sending build requests to the worker pool
var BuildSendWorkerChan = make(chan *models.Buildrequest, BuildWorkerPool)

// BuildRecieveWorkerChan channel is passed to the worker, and will recieve work to do
var BuildRecieveWorkerChan chan chan *models.Buildrequest

var DockerCleanSendWorkerChan = make(chan *DockerClean, BuildWorkerPool)
var DockerCleanRecieveWorkerChan chan chan *DockerClean

//StartDispatcher will create a number of workers and listen on a channel for new work to send to workers
func StartDispatcher(count int) {

	// initialize the channel we are going to put the workers' work channels into.
	BuildRecieveWorkerChan = make(chan chan *models.Buildrequest, count)
	DockerCleanRecieveWorkerChan = make(chan chan *DockerClean, count)

	for i := 1; i < count+1; i++ {
		log.Debugf("Dispatcher starting worker%d", i)
		worker := NewCommandWorker(i, BuildRecieveWorkerChan)
		worker.Start()
	}

	// start a single docker worker
	log.Debugf("Dispatcher starting worker%d", 0)
	worker := NewDockerWorker(0, DockerCleanRecieveWorkerChan)
	worker.Start()

	go func() {
		// forever
		for {
			// recieve build requests and docker clean requests, notify
			// a worker via it's work channel, and repeat
			select {
			// get work sent on the worker channel
			case buildrequest := <-BuildSendWorkerChan:

				log.Debugf("Received build request", BuildRequestToString(*buildrequest))
				go func() {
					// get a worker's channel
					worker := <-BuildRecieveWorkerChan

					// send the work
					worker <- buildrequest

				}()

			// get work sent on the docker clean channel
			case cleanrequest := <-DockerCleanSendWorkerChan:
				log.Debugf("Received docker clean requeust %v", cleanrequest)
				go func() {
					worker := <-DockerCleanRecieveWorkerChan
					worker <- cleanrequest
				}()

			}
		}
	}()
}
