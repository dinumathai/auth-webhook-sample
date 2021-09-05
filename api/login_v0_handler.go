package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dinumathai/auth-webhook-sample/auth"
	"github.com/dinumathai/auth-webhook-sample/log"
	"github.com/dinumathai/auth-webhook-sample/types"
	"github.com/dinumathai/auth-webhook-sample/util/response"
	"gopkg.in/yaml.v2"
)

var userConfig map[string]types.UserDetails

// LoginV0Handler -- Handle auth using property file. For testing only
func LoginV0Handler(config *types.ConfigMap) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		defer log.Debugf("LoginV0Handler Elapsed - %s", time.Since(start))

		//Check for valid username and password
		username, password, ok := r.BasicAuth()
		if !ok {
			sendResponse(http.StatusUnauthorized, "", types.RawAuthResponse{}, fmt.Errorf("Need valid username and password as basic auth"), w)
			return
		}
		userDetailFromConfig, err := validateAndGetUser(config, username, password)
		if err != nil {
			errHandle(w, fmt.Sprintf("Unable to validate : %s", err), "Authentication failed", 401)
			return
		}
		user := types.User{
			Username: userDetailFromConfig.UserName,
			EMail:    userDetailFromConfig.Email,
			UID:      userDetailFromConfig.UID,
			Groups:   userDetailFromConfig.Groups}
		if user.UID == "" {
			user.UID = user.Username
		}
		token, err := auth.GenerateToken(user, "", auth.V0)
		if err != nil {
			errHandle(w, fmt.Sprintf("Something is wrong with auth token. : %s", err), "Authentication failed", 401)
			return
		}

		v1Token := types.V1Token{
			Token:  token.JWT,
			Expiry: token.Expiry,
		}

		data, _ := json.Marshal(v1Token)
		response := JSONResponse{}
		response.status = http.StatusCreated
		response.data = data

		response.Write(w)
	}
}

func validateAndGetUser(config *types.ConfigMap, userName, password string) (types.UserDetails, error) {
	userDtlMap, err := getV0UserConfig(config)
	if err != nil {
		return types.UserDetails{}, err
	}
	if userDtl, ok := userDtlMap[userName]; ok {
		if userDtl.Password == password {
			if userDtl.UserName == "" {
				userDtl.UserName = userName
			}
			return userDtl, nil
		} else {
			return types.UserDetails{}, errors.New("Invalid Credentials")
		}
	}
	return types.UserDetails{}, errors.New("User Not present")
}

func getV0UserConfig(config *types.ConfigMap) (map[string]types.UserDetails, error) {
	if config.AuthConfig.V0.UserDetailFilePath == "" {
		return userConfig, errors.New("Invalid Config")
	}
	//TODO : Avoid subsequent read if first time read/parsing failed
	if userConfig != nil {
		return userConfig, nil
	}
	data, err := ioutil.ReadFile(config.AuthConfig.V0.UserDetailFilePath)
	if err != nil {
		log.Errorf("User Details config read Failed: %v", err)
		return userConfig, err
	}
	var userConf types.UserDetailsConfig
	if yamlErr := yaml.Unmarshal(data, &userConf); yamlErr != nil {
		log.Errorf("Error deserializing yaml %v", yamlErr)
		return nil, yamlErr
	}
	userConfig = userConf.UserDetails
	return userConfig, nil
}

// errHandle packages an error into an http response
func errHandle(w http.ResponseWriter, longmsg string, shortmsg string, status int) {
	log.Errorf(longmsg)
	errorResponse := ErrorResponse{
		Status:       status,
		ErrorMessage: shortmsg,
	}
	data, _ := json.Marshal(errorResponse)
	response := JSONResponse{}
	response.status = http.StatusUnauthorized
	response.data = data
	response.Write(w)
}

func sendResponse(statusCode int, provider string, rawAuthResponse types.RawAuthResponse, err error, w http.ResponseWriter) {
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
		response := rawAuthResponse.Provider.Token
		resData, _ = json.Marshal(response)
	}

	res.Data = resData
	res.Write(w)
}
