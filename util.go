package packerd

import (
	"encoding/json"
	"fmt"

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
