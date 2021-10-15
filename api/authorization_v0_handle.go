package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dinumathai/auth-webhook-sample/auth"
	"github.com/dinumathai/auth-webhook-sample/log"
	"github.com/dinumathai/auth-webhook-sample/types"
)

// AuthorizeV0Handler -- Handle authentication using property file. For testing only
func AuthorizeV0Handler(apiVersion auth.Version) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer log.Debugf("AuthorizeV0Handler Elapsed - %s", time.Since(start))

		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Debugf("Error in Read of request body : %s", err)
			sentAuthorizationResponse(w, false, "Error in Read of request body")
			return
		}
		rawContent := json.RawMessage(string(content))
		log.Debugf("Request body : %s", rawContent)
		log.Debugf("Request headers : %v", r.Header)

		sentAuthorizationResponse(w, true, "")
	}
}

func sentAuthorizationResponse(w http.ResponseWriter, allowed bool, reason string) {
	response := types.AuthorizationResponse{
		APIVersion: "authorization.k8s.io/v1",
		Kind:       "SubjectAccessReview",
	}
	if allowed {
		response.Status = &types.AuthorizationStatus{
			Allowed: true,
		}
	} else {
		response.Status = &types.AuthorizationStatus{
			Allowed: false,
			Denied:  true,
			Reason:  reason,
		}
	}
	responseBytes, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)

	w.Write(responseBytes)
}
