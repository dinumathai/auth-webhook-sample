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

//LDAP service account for LDAP configured as V1
type LDAPServiceAccount struct {
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
}

type LdapConfig struct {
	LdapUserFilter  string   `yaml:"ldapUserFilter"`
	LdapGroupFilter string   `yaml:"ldapGroupFilter"`
	SvcAccPrefix    string   `yaml:"svcAccPrefix"`
	LdapAttributes  []string `yaml:"ldapAttributes"`
	LdapServiceOU   string   `yaml:"ldapServiceOU"`
	LdapOU          string   `yaml:"ldapOU"`
	LdapProdOU      string   `yaml:"ldapProdOU"`
	LdapV1UserOU    string   `yaml:"ldapV1UserOU"`
	LdapV2UserOU    string   `yaml:"ldapV2UserOU"`
	ProdAccKey      string   `yaml:"prodAccKey"`
	LdapAccKey      string   `yaml:"ldapAccKey"`
}

//AdMeta ...
type AdMeta struct {
	LDAPHost     string `yaml:"ldapHost"`
	LDAPPort     int    `yaml:"ldapPort"`
	UseTLS       bool   `yaml:"useTLS"`
	InsecureSkip bool   `yaml:"insecureSkip"`
	Provider     string `yaml:"provider"`
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

// V2LoginResponseJSON is the structure returned by a V2 Login endpoint
type V2LoginResponseJSON struct {
	Providers []Provider `json:"provider,omitempty"` // Note mixed singular of Provider
}

//V1Token ...
type V1Token struct {
	Token  string `json:"token,omitempty"`
	Expiry int64  `json:"expiry,omitempty"`
}

//RequestParameters options for this service
type RequestParameters struct {
	Provider  string
	ClaimType string
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

// RoleName are just strings
//type RoleName string

// Roles are simplified to just role: resource1...resourceN style
//type Roles security.Roles //map[RoleName][]string

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
	Iat      int64               `json:"iat"`
	UID      string              `json:"uid"`
	Username string              `json:"username"`
	Expiry   int64               `json:"exp"`
	Groups   []string            `json:"groups"`
	Apps     []AppOwner          `json:"apps,omitempty"`
	Roles    map[string][]string `json:"roles,omitempty"`
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

// AppOwner says which user owns each app
type AppOwner struct {
	App   string `json:"app,omitempty"`
	Owner string `json:"owner,omitempty"`
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

//UserStore ..
type UserStore struct {
	Provider string      `json:"provider,omitempty"`
	MetaUser interface{} `json:"metaUser,omitempty"`
}

// Protocol is the type of auth protocol
type Protocol int

// Possible values of Protocol
const (
	Unknown   Protocol = iota
	LDAP               // something via LDAP
	Radius             // Direct radius
	ActiveDir          // Direct AD connection
	OSUser             // OS provided
	FileUP             // Username and password from local file
)

// ProviderID is the type of the unique id of a provider
type ProviderID string

// Provider2 type
type Provider2 struct {
	ID          ProviderID `json:"id" validate:"nonzero"`       // unique short identifer for this provider, machine and UI use
	Description string     `json:"description"`                 // descriptive text about this provider, UI use only
	Protocol    Protocol   `json:"protocol" validate:"nonzero"` // protocol this provider speaks
	Address     string     `json:"address"`                     // url/IP & port or filename
	Priority    int        `json:"priority"`                    // sequence of this provider, higher ones should be shown first
	PolicyID    string     `json:"policyid"`                    // policy id of the policy to apply to users who arrive via this provider
}
