package config

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/dinumathai/auth-webhook-sample/log"
	"github.com/dinumathai/auth-webhook-sample/types"

	"github.com/ghodss/yaml"
)

var AppConfig = types.ConfigMap{}

// Load will load configuration from k8s config map.
// If it did not find any then it will fall back to config map provided with binary
func Load() (*types.ConfigMap, error) {
	var config types.ConfigMap

	data, err := ReadConfigData("Main auth data", "CONFIG_FILE", "auth_config.yaml")
	if err != nil {
		return &config, err
	}
	if yamlErr := yaml.Unmarshal(data, &config); yamlErr != nil {
		log.Fatalf("Error deserializing yaml data starting %s: %s", string(data[0:20]), yamlErr.Error())
		return nil, yamlErr
	}

	signingKey := os.Getenv("AUTH_SIGING_KEY")
	if len(signingKey) != 0 {
		config.AuthConfig.AuthSigningKey = signingKey
	}
	if len(config.AuthConfig.AuthSigningKey) == 0 {
		log.Fatal("Unable to get SigningKey from environment variable(AUTH_SIGING_KEY) or application configuration")
	}
	AppConfig = config
	return &config, nil
}

// ReadConfigData will read the data from a k8s config map.
//   If it did not find any then it will fall back to a file provided with binary.
//   Name: descriptive string for log messages
//   Envvar: the name of the environment var where the relative path to the config file is stored
func ReadConfigData(name string, envvar string, fallback string) ([]byte, error) {

	var data []byte

	fileP := os.Getenv(envvar) //Get File path from Env Var

	if fileP == "" { // None derived from env var & mapped k8s config-map so fall back
		wd, _ := os.Getwd()
		fileP = filepath.Join(wd, "config", fallback)
		log.Infof("No k8s %s given in %s, falling back to: %s", name, envvar, fileP)
	} else {
		fileP = filepath.FromSlash(fileP)
		log.Infof("Loading %s from k8s config Map: %s", name, fileP)
		name += " via k8s"
	}

	data, err := ioutil.ReadFile(fileP)
	if err != nil {
		log.Fatalf("Error reading %s from %s: %s", name, fileP, err.Error())
		return data, err
	}

	return data, nil

}
