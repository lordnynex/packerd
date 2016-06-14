package restapi

import (
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/go-openapi/runtime/middleware"
	"github.com/tompscanlan/packerd"
	"github.com/tompscanlan/packerd/models"
	"github.com/tompscanlan/packerd/restapi/operations/informational"
)

func health(params informational.GetHealthParams) middleware.Responder {
	h := GetHealth()

	health := informational.NewGetHealthOK()
	health.SetPayload(h)
	log.Debugf("%v", health.Payload)

	return health
}

// GetHealth returns the current health status of the packerd service
func GetHealth() *models.Health {
	var health *models.Health

	health = new(models.Health)
	health.Diskpercentfull = intptr(100)
	health.Status = strptr("")

	err, percentfull := packerd.GetDirPercentFull(packerd.DockerDir)
	if err != nil {
		health.Diskpercentfull = intptr(100)
		health.Status = strptr("err")
		return health
	}

	health.Goroutines = uint64(runtime.NumGoroutine())

	*health.Diskpercentfull = percentfull

	switch {
	case *health.Diskpercentfull > 95:
		*health.Status = "bad"
	case *health.Diskpercentfull > 90:
		*health.Status = "poor"
	default:
		*health.Status = "ok"
	}

	return health
}

func intptr(i int) *int64 {
	conv := int64(i)
	return &conv
}
func strptr(s string) *string {
	return &s
}
