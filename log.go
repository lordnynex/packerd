package packerd

import (
	"log"
	"os"
)

// global logger
var Logger = log.New(os.Stderr, "debug: ", log.Ldate|log.Ltime|log.Lshortfile)
