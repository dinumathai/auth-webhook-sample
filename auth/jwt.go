package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/dinumathai/auth-webhook-sample/config"
	"github.com/dinumathai/auth-webhook-sample/log"
	"github.com/dinumathai/auth-webhook-sample/types"

	jwt "github.com/dgrijalva/jwt-go"
)

//BearerSchema that this service is expecting
const (
	BearerSchema string = "Bearer "
	SigningKey          = "AUTH_SIGNING_KEY"

	// JWT/API Version constants
	V0 = 0
	V1 = 1
	V2 = 2

	APIVerString = "authentication.k8s.io/v2"
)

// Version -- constrained type
type Version int

//GenerateToken generates a full JWT groups and apps etc.
func GenerateToken(user types.User, hclaims string, majVersion Version) (types.Token, error) {

	//Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	// Set claims
	claims["username"] = user.Username
	claims["uid"] = user.UID
	//	filteredGroups := FilterGroupsOnClaims(user.Groups, hclaims)
	claims["groups"] = user.Groups
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iat"] = time.Now().Unix()

	signedToken, err := token.SignedString([]byte(config.AppConfig.AuthConfig.AuthSigningKey))
	if err != nil {
		log.Errorf("Cannot sign token  : %s", err)
		return types.Token{}, err
	}

	return types.Token{
		JWT:    signedToken,
		Expiry: claims["exp"].(int64),
	}, nil

}

//ValidateToken validates JWT token provided by user and fills out the UserInfo structure from the data within
func ValidateToken(req *http.Request, majVersion Version) (types.UserInfo, int, error) {
	auth := new(bool)
	*auth = false

	errUserInfo, errBadReq := types.UserInfo{
		APIVersion: APIVerString,
		Kind:       "TokenReview",
		Status: &types.Status{
			Authenticated: auth,
			User:          nil,
		},
	}, fmt.Errorf("Either need valid JWT bearer token in Authorization header or need valid kubernetes webhook auth request (Please refer - %s)", "https://kubernetes.io/docs/reference/access-authn-authz/authentication/#webhook-token-authentication")

	//If body is empty or not able to parse properly then try with Auth header
	request, err := getRequestBody(req.Body)
	if err != nil {
		log.Debugf("Unable to parse request body: %v. Trying with Authorization header.", err)

		token, err := checkAuthScheme(req.Header.Get("Authorization"))
		if err != nil {
			return errUserInfo, http.StatusBadRequest, errBadReq
		}

		return validate(token, majVersion) // note: most work happens here <<<
	}

	log.Debug("Received auth token from body. Skipping Auth Header check.")
	//Get Auth token from body and validate
	if request.Spec.Token != "" {
		return validate(request.Spec.Token, majVersion) // note: most work happens here <<<
	}
	return errUserInfo, http.StatusBadRequest, errBadReq
}

func getRequestBody(body io.ReadCloser) (types.Request, error) {
	content, err := ioutil.ReadAll(body)
	if err != nil {
		log.Debugf("Error in Read of request body : %s", err)
		return types.Request{}, err
	}
	rawContent := json.RawMessage(string(content))
	log.Debugf("Request body : %s", rawContent)
	marshaledContent, err := rawContent.MarshalJSON()
	if err != nil {
		log.Debugf("Error in marshaling request body : %s", err)
		log.Debugf("Request Body might be empty. If so we will try with Authorization Header")
		return types.Request{}, err
	}

	var request types.Request
	err = json.Unmarshal(marshaledContent, &request)
	if err != nil {
		log.Debugf("Error in un-marshaling request body : %s", err)
		log.Debugf("Request Body might be empty or not of kube webhook auth request type. If so we will try with Authorization Header")
		return types.Request{}, err
	}

	return request, nil
}

// validate does much of the work of ValidateToken
func validate(bearerToken string, majVersion Version) (types.UserInfo, int, error) {

	signingKey := config.AppConfig.AuthConfig.AuthSigningKey

	var auth bool
	var claims types.JWTClaimsJSON // special struct for decoding the json

	// user we'll return, initially in error state
	u := types.UserInfo{
		APIVersion: "authentication.k8s.io/v1beta1",
		Kind:       "TokenReview",
		Status: &types.Status{
			Authenticated: &auth,
			User:          nil,
		},
	}

	token, err := jwt.ParseWithClaims(bearerToken, &claims, func(token *jwt.Token) (interface{}, error) {
		if !strings.HasPrefix(token.Method.Alg(), "HS") { // HMAC are the only allowed signing methods
			log.Errorf("Unexpected signing method: %s", token.Method.Alg())
			return nil, fmt.Errorf("Unexpected signing method: %s", token.Method.Alg())
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		log.Errorf("Error Parsing JWT. Error - %v", err)
		return u, http.StatusBadRequest, err
	}

	if !token.Valid {
		log.Errorf("Token not valid: %v", err)
		return u, http.StatusBadRequest, err
	}

	// Token is valid so fill in the rest of u with happy state and return it
	auth = true
	u.Status.Authenticated = &auth
	u.Status.User = &types.User{
		Username: claims.Username,
		UID:      claims.UID,
		Groups:   claims.Groups}

	return u, http.StatusOK, nil

}

func checkAuthScheme(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("Didn't receive any auth token")
	}

	// Confirm the request is sending Bearer Authentication credentials.
	if !strings.HasPrefix(authHeader, BearerSchema) {
		return "", errors.New("Authorization requires 'Bearer' scheme")
	}

	// Get the token from the request header
	// The first six characters are skipped - e.g. "Bearer ".
	return authHeader[len(BearerSchema):], nil
}

// FilterGroupsOnClaims returns groups that match at least one claim
func FilterGroupsOnClaims(groups []string, claims string) []string {

	c2 := strings.TrimSpace(claims)
	if len(c2) == 0 { // no claims -> everything
		return groups
	}

	clist := strings.Split(c2, ",")

	gs := make([]string, 0)
grouploop:
	for _, g := range groups {
		for _, c := range clist {
			cmi := strings.TrimSpace(c)
			if len(cmi) > 0 && strings.Contains(g, cmi) {
				gs = append(gs, g)
				continue grouploop
			}
		}
	}

	return gs

}
