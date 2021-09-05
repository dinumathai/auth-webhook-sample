package response

import (
	"encoding/json"

	"net/http"

	"github.com/dinumathai/auth-webhook-sample/log"
)

//Writer interface
type Writer interface {
	Write(w http.ResponseWriter)
}

//Response defines a response sent to client
type Response struct {
	Status      int
	ContentType string
	Data        []byte
}

// var log = logger.GetLogger()

// SendJSON the HTTP response
func SendJSON(statusCode int, result interface{}, w http.ResponseWriter) {
	responseData, _ := json.Marshal(result)
	Send(statusCode, nil, responseData, w)
}

// Send the HTTP response
func Send(statusCode int, err error, result []byte, w http.ResponseWriter) {
	// re := string(result)
	// log.Print(re)
	var responseData []byte
	if err != nil {
		res := ErrorResponse{
			ErrorMessage: err.Error(),
		}
		responseData, _ = json.Marshal(res)
	} else {
		responseData = result
	}

	res := Response{
		Status: statusCode,
		Data:   responseData,
	}
	res.Write(w)
}

//ErrorResponse defines error if in case any
type ErrorResponse struct {
	Status       int    `json:"status,omitempty"`
	ErrorMessage string `json:"error,omitempty"`
}

func (jr *Response) Write(w http.ResponseWriter) {
	// defer log.Infof("HTTP Status : %v", jr.Status)
	log.Debugf("HTTP Status : %v", jr.Status)

	if jr.ContentType == "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	} else {
		w.Header().Set("Content-Type", jr.ContentType)
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(jr.Status)

	w.Write(jr.Data)
}
