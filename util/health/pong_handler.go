package health

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dinumathai/auth-webhook-sample/log"
	resp "github.com/dinumathai/auth-webhook-sample/util/response"
)

//Response ...
type Response struct {
	BuildVersion string `json:"build.version"`
}

var (
	//Version ...
	Version = "NotSet"
)

//PongHandler checks health of service
func PongHandler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	healthResponse := &Response{
		BuildVersion: Version,
	}

	data, _ := json.Marshal(healthResponse)

	response := resp.Response{}
	response.Status = http.StatusOK
	response.Data = data

	response.Write(w)
	log.Debugf("Responded to health check! %s %s %s", r.Method, r.URL.String(), time.Since(start).String())
}
