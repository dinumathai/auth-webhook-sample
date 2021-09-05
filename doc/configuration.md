## Configuration

## Configuration And Running the application
The service reads the configuration file name from the environment variable `CONFIG_FILE`. If the environment variable is not set the default value is `$PWD/config/auth_config.yaml`.


To run the application set the `CONFIG_FILE` and run the executable.
```
export CONFIG_FILE=/the/path/to/auth_config.yaml
./auth-webhook-sample
```

## Run in https mode
The auth service expects the certificate and key in environment variables `AUTH_CERT_TLS_CRT` and `AUTH_CERT_TLS_KEY` respectively. If the environment variables are set the service will start in `https` mode instead of `http`. While deploying in Kubernetes its preferred to store the certificates as Kubernetes secrets and make it available to the container as above environment variables.


## Config file
The service reads the configuration file name from the environment variable `CONFIG_FILE`. If the environment variable is not set the default value is `$PWD/config/auth_config.yaml`.

Sample config available at [config/auth_config.yaml](../config/auth_config.yaml)

| Configuration | Type | Mandatory | Description |
| ------------  | ---- | --------- | ----------  |
| authConfig.serverAddress | int | Mandatory | The port number in which the application is going to listen. |
| authConfig.v0.userDetailFilePath | string | Mandatory | For V0 api - The path of the file that holds user details. Refer [config/user_details.yaml](../config/user_details.yaml)|
| authConfig.authSigningKey | string | Mandatory | The Signing Key for generating the auth token. |