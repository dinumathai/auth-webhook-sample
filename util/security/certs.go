package security

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"

	"github.com/dinumathai/auth-webhook-sample/log"
)

var basePath string

// CrtPath will be the path to the cert file for tls
var CrtPath string

// KeyPath will be the path to the cert file for tls
var KeyPath string

func init() {
	basePath = getBasePath()
	CrtPath = basePath + string(os.PathSeparator) + "tls.crt"
	KeyPath = basePath + string(os.PathSeparator) + "tls.key"
}

//LoadCertsFromEnvVariable - Pull and Load SSL Certs for HTTPS..
func LoadCertsFromEnvVariable(crtEnvVar string, keyEnvVar string) error {
	err := createDirectoryIfNotExist(getBasePath())
	if err != nil {
		return err
	}

	err = writeFile(GetCertificatePath(), crtEnvVar)
	if err != nil {
		return err
	}

	err = writeFile(GetKeyPath(), keyEnvVar)
	if err != nil {
		return err
	}

	return nil
}

func createDirectoryIfNotExist(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeFile(path string, envVar string) error {
	decodedString, err := decodeIfEncoded(os.Getenv(envVar))
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, []byte(decodedString), 0644)

	return err
}

func decodeIfEncoded(envVar string) (string, error) {
	match, err := regexp.MatchString("^([A-Za-z0-9+/]{4})*([A-Za-z0-9+/]{4}|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{2}==)$", envVar) // Regex to check if string in base 64 encoded or not
	if err != nil {
		log.Errorf("Something went wrong while mathching RegEx : %v", err)
		return "", err
	}
	if match {
		sDec, decodeErr := b64.StdEncoding.DecodeString(envVar)
		if decodeErr != nil {
			log.Errorf("Something went wrong while decoding string : %v", decodeErr)
			return "", decodeErr
		}
		return string(sDec), nil
	}
	return envVar, nil
}

func getBasePath() string {
	if runtime.GOOS == "windows" {
		return "C:\\cert-temp-auth" // kinda lame
	} else {
		return "/tmp/cert-temp-auth"
	}
}

func GetCertificatePath() string {
	return getBasePath() + string(os.PathSeparator) + "tls.crt"
}

func GetKeyPath() string {
	return getBasePath() + string(os.PathSeparator) + "tls.key"
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

// GenerateCertificate for a given template
func GenerateCertificate(cert x509.Certificate) (string, string, error) {
	// priv, err := rsa.GenerateKey(rand.Reader, *rsaBits)
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil {
		log.Error(err)
		return "", "", err
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &cert, &cert, publicKey(priv), priv)
	if err != nil {
		log.Errorf("Failed to create certificate: %s", err)
		return "", "", err
	}
	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	crt := out.String()
	out.Reset()
	pem.Encode(out, pemBlockForKey(priv))
	key := out.String()
	return crt, key, nil
}
