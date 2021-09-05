package server

import (
	"net/http"
	"os"
	"strconv"

	"github.com/dinumathai/auth-webhook-sample/types"
	"github.com/dinumathai/auth-webhook-sample/util/routing"
	"github.com/dinumathai/auth-webhook-sample/util/security"

	"github.com/dinumathai/auth-webhook-sample/log"
)

const (
	authSSLCrtEnvVar = "AUTH_CERT_TLS_CRT"
	authSSLKeyEnvVar = "AUTH_CERT_TLS_KEY"
)

//Start starts the server
func Start(config *types.ConfigMap) {
	if (os.Getenv(authSSLCrtEnvVar) != "") && (os.Getenv(authSSLKeyEnvVar) != "") {
		log.Infof("Loading HTTPS certificates....")
		err := security.LoadCertsFromEnvVariable(authSSLCrtEnvVar, authSSLKeyEnvVar)
		if err != nil {
			log.Infof("SSL Certs Not available...")
		} else {
			log.Infof("SSL Certs loaded successfully...")
		}
	}

	router := routing.BuildRouter(BuildRoutes(config))
	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./swaggerui/"))))
	http.Handle("/", router)

	log.Infof("Starting Server...")
	var err error
	if (os.Getenv(authSSLCrtEnvVar) != "") && (os.Getenv(authSSLKeyEnvVar) != "") {
		log.Info("Starting server with SSL on port ", config.AuthConfig.ServerAddress)
		err = http.ListenAndServeTLS(":"+strconv.Itoa(config.AuthConfig.ServerAddress), security.CrtPath, security.KeyPath, nil)
	} else {
		log.Info("DEV MODE - Starting HTTP server on port ", config.AuthConfig.ServerAddress)
		err = http.ListenAndServe(":"+strconv.Itoa(config.AuthConfig.ServerAddress), nil)
	}
	if err != nil {
		log.Info("Starting server - Failed : " + err.Error())
	}
}
