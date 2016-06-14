package packerd

import (
	log "github.com/Sirupsen/logrus"

	"time"
)

type DockerClean struct {
	Ticker *time.Ticker
}

func init() {

	DockerClean := NewDockerClean()
	go func() {
		for {

			select {
			// when the ticker pings us, send a clean job to the workers
			case work := <-DockerClean.Ticker.C:
				DockerCleanSendWorkerChan <- DockerClean
				log.Info("in docker clean, got the ticker, ", work)
			}
		}
	}()
}
func NewDockerClean() *DockerClean {
	d := new(DockerClean)
	d.Ticker = time.NewTicker(DockerCleanInterval)

	return d
}
