package api

import (
	"encoding/json"
	"net/http"

	"github.com/dinumathai/auth-webhook-sample/auth"
	"github.com/dinumathai/auth-webhook-sample/types"
	"github.com/dinumathai/auth-webhook-sample/util/response"
)

//ValidationHandler validates the token
func ValidationHandler(apiVersion auth.Version) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userInfo, statusCode, tokenErr := auth.ValidateToken(r, apiVersion)
		sendV1BetaResponse(statusCode, userInfo, tokenErr, w)
	}
}

func sendV1BetaResponse(statusCode int, userInfo types.UserInfo, err error, w http.ResponseWriter) {
	res := response.Response{}

	switch statusCode {
	case 200:
		res.Status = http.StatusOK
	case 201:
		res.Status = http.StatusCreated
	case 400:
		res.Status = http.StatusBadRequest
	case 401:
		res.Status = http.StatusUnauthorized
	case 403:
		res.Status = http.StatusForbidden
	case 404:
		res.Status = http.StatusNotFound
	case 500:
		res.Status = http.StatusInternalServerError
	}

	var resData []byte
	if err != nil {
		resData, _ = json.Marshal(types.AuthResponse{
			Error: err.Error(),
		})
	} else {
		resData, _ = json.Marshal(userInfo)
	}

	res.Data = resData
	res.Write(w)
}
