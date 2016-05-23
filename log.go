package packerd

import (
	"os"

	log "github.com/Sirupsen/logrus"
)

func init() {
	// global logger
	//var Logger = log.New(os.Stderr, "debug: ", log.Ldate|log.Ltime|log.Lshortfile)
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stderr)
	//log.SetLevel(log.WarnLevel)
	log.SetLevel(log.InfoLevel)
}
