package main

import (
	"flag"

	cfg "github.com/dinumathai/auth-webhook-sample/config"
	"github.com/dinumathai/auth-webhook-sample/log"
	"github.com/dinumathai/auth-webhook-sample/server"
	"github.com/dinumathai/auth-webhook-sample/util/health"
)

var addr = flag.String("listen-address", ":8443", "The address to listen on for HTTPS requests.")

// Version is the program build version set in Dockerfile via ldflags
var Version = "notSet"

func main() {
	//Get config
	config, err := cfg.Load()
	if err != nil || config.AuthConfig.ServerAddress == 0 {
		log.Fatalf("Config not loaded correctly - %v", err)
	}

	//server
	log.Info("Starting Auth server..........")
	health.Version = Version
	// server.Start(userStore, config) //Commeting out local cache for now.
	server.Start(config)
}
