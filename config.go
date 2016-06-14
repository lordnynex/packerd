package packerd

import (
	"time"

	log "github.com/Sirupsen/logrus"
)

// provide a nice way to convert to various orders of magnitude
const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// DockerDir is where docker's data resides, used for health checking based on space available
var DockerDir = "/var/lib/docker"

// HealthyFreeDiskGB is how much space should be available to be healthy
var HealthyFreeDiskGB = 300

var DefaultLogLevel = log.DebugLevel

//var DefaultLogLevel = log.InfoLevel

// how big a backlog to allow for builds
var BuildWorkerPool = 100

// number of workers
var WorkerCount = 5

var DockerCleanInterval = 15 * time.Second
