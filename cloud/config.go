package cloud

import (
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	MqttBrokerURL   = "tls://mqtt.googleapis.com:443"
	CACertsURL      = "https://pki.goog/roots.pem"
	ProtocolVersion = 4 // corresponds to MQTT 3.1.1
	ReconnectDelay  = 20 * time.Second
	TokenDuration   = 24 * time.Hour // maximum 24 h
)

var (
	once       sync.Once
	caCertPool *x509.CertPool
)

// Config stores the data for connecting to the MQTT broker.
type Config struct {
	ProjectID  string `yaml:"projectID"`
	Region     string `yaml:"region"`
	RegistryID string `yaml:"registryID"`
	DeviceID   string `yaml:"deviceID"`
	PrivateKey string `yaml:"privateKey"`
	Algorithm  string `yaml:"algorithm"`
	StorePath  string `yaml:"storePath"`
}

// CreateJWT creates a Cloud IoT Core JWT for the given project id.
func CreateJWT(projectID, privateKey, algorithm string, iat, exp time.Time) (string, error) {

	claims := jwt.StandardClaims{
		Audience:  projectID,
		IssuedAt:  iat.Unix(),
		ExpiresAt: exp.Unix(),
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod(algorithm), claims)

	switch algorithm {
	case "RS256":
		privKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
		if err != nil {
			return "", err
		}
		return token.SignedString(privKey)
	case "ES256":
		privKey, err := jwt.ParseECPrivateKeyFromPEM([]byte(privateKey))
		if err != nil {
			return "", err
		}
		return token.SignedString(privKey)
	}
	return "", errors.New("invalid token algorithm, need one of ES256 or RS256")
}

// GetCACertPool gets the latest Google CA certs from 'https://pki.goog/roots.pem'
func GetCACertPool() *x509.CertPool {
	once.Do(func() {
		// Load Google Root CA certs.
		resp, err := http.Get(CACertsURL)
		if err != nil {
			log.Println(errors.Wrapf(err, "cannot get cert file from '%s'", CACertsURL))
			return
		}
		defer resp.Body.Close()
		certs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(errors.Wrapf(err, "cannot read cert file from '%s'", CACertsURL))
			return
		}
		roots := x509.NewCertPool()
		if ok := roots.AppendCertsFromPEM(certs); !ok {
			log.Println(errors.Errorf("cannot parse cert file from '%s'", CACertsURL))
			return
		}
		caCertPool = roots
	})

	return caCertPool
}
