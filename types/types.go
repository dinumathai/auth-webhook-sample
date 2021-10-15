package types

import (
	"fmt"
	"time"
)

//ConfigMap maps all the configuration needed by auth
type ConfigMap struct {
	AuthConfig AuthConfig `yaml:"authConfig"`
}

//AuthConfig ...
type AuthConfig struct {
	V0             UserMeta `yaml:"v0"`
	ServerAddress  int      `yaml:"serverAddress"`
	AuthSigningKey string   `yaml:"authSigningKey"`
}

// UserMeta - User detail for V0 api
type UserMeta struct {
	Source             string `yaml:"source"`
	UserDetailFilePath string `yaml:"userDetailFilePath"`
}

//AuthResponse ...
type AuthResponse struct {
	Provider Provider `json:"provider,omitempty"`
	Error    string   `json:"error,omitempty"`
}

//Provider ...
type Provider struct {
	Name  string `json:"name,omitempty"`
	Token Token  `json:"token,omitempty"`
}

//Token return JWT with its expiry time
type Token struct {
	JWT    string `json:"jwt,omitempty"`
	Expiry int64  `json:"expiry,omitempty"`
}

//V1Token ...
type V1Token struct {
	Token  string `json:"token,omitempty"`
	Expiry int64  `json:"expiry,omitempty"`
}

//RawAuthResponse ...
type RawAuthResponse struct {
	Provider   Provider
	Error      error
	HTTPStatus int
}

//UserInfo holds authentication information
type UserInfo struct {
	APIVersion string  `json:"apiVersion,omitempty"`
	Kind       string  `json:"kind,omitempty"`
	Status     *Status `json:"status,omitempty"`
}

//User holds user information from AD
type User struct {
	Username string   `json:"username,omitempty"`
	EMail    string   `json:"email,omitempty"`
	UID      string   `json:"uid,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

// UserCredentials is a simple username/password pair
type UserCredentials struct {
	UserName string
	Password string
}

// JWTClaimsJSON is used for decoding an incoming JSON JWT payload to the /authenticate API
type JWTClaimsJSON struct {
	Iat      int64    `json:"iat"`
	UID      string   `json:"uid"`
	Username string   `json:"username"`
	Expiry   int64    `json:"exp"`
	Groups   []string `json:"groups"`
}

// Valid so that JWTClaimsJSON satisfies the jwt.Claims interface
func (c JWTClaimsJSON) Valid() error {
	if c.UID == "" {
		return fmt.Errorf("UID must be present in token claims")
	}
	if c.Expiry == 0 {
		return fmt.Errorf("Token has no expiry")
	}
	if c.Expiry < int64(time.Now().Unix()) {
		return fmt.Errorf("Token has expired")
	}
	if c.Iat > int64(time.Now().Unix()+int64(time.Second)) {
		return fmt.Errorf("Token is from the future")
	}
	return nil
}

//Status indicates if user is authenticated or not
type Status struct {
	Authenticated *bool `json:"authenticated,omitempty"`
	User          *User `json:"user,omitempty"`
}

//Request maps the incoming auth request from api-server
type Request struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Spec       *Spec  `json:"spec,omitempty"`
}

//Spec maps to the bearer token send by api-server
type Spec struct {
	Token string `json:"token,omitempty"`
}

//Authorization response
type AuthorizationResponse struct {
	APIVersion string               `json:"apiVersion,omitempty"`
	Kind       string               `json:"kind,omitempty"`
	Status     *AuthorizationStatus `json:"status,omitempty"`
}

type AuthorizationStatus struct {
	Allowed bool   `json:"allowed,omitempty"`
	Denied  bool   `json:"denied,omitempty"`
	Reason  string `json:"reason,omitempty"`
}
