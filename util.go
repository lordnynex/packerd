package packerd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"runtime"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/tompscanlan/packerd/models"
)

func BuildRequestToString(br models.Buildrequest) string {

	b, err := json.Marshal(br)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal to json: %s", err.Error())
		return msg
	}

	return string(b[:])
}

func StreamToLog(reader io.Reader) {
	b := bufio.NewScanner(reader)
	for b.Scan() {
		log.WithFields(log.Fields{
			"function": "StreamToLog",
		}).Infoln(b.Text())
	}

	if err := b.Err(); err != nil {
		log.WithFields(log.Fields{
			"function": "StreamToLog",
		}).Errorln("error reading:", err)
	}
}

func StreamToString(reader io.Reader, s *string) {
	var mutex = &sync.Mutex{}

	b := bufio.NewScanner(reader)
	for b.Scan() {
		mutex.Lock()
		*s = *s + b.Text()
		mutex.Unlock()
		runtime.Gosched()
	}

	if err := b.Err(); err != nil {
		*s = *s + fmt.Sprintf("%v", err)
	}
}
