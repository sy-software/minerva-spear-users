package domain

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"os"

	"github.com/rs/zerolog/log"
)

type Token struct {
	// Token duration in seconds, default: 7 days
	Duration int64 `json:"duration,omitempty"`
	// Refresh Token in seconds, default: 30 days
	RefreshDuration int64 `json:"refreshDuration,omitempty"`
	// For JWT signature using RS256 algorithm
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`

	// We need to parse the key into *rsa.PrivateKey to be usable
	rsaKey *rsa.PrivateKey
}

// KeyPair parses private and public key string into a *rsa.PrivateKey instance
func (t *Token) KeyPair() (*rsa.PrivateKey, error) {
	if t.rsaKey != nil {
		return t.rsaKey, nil
	}

	privPem, _ := pem.Decode([]byte(t.PrivateKey))

	if privPem.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("RSA private key is of the wrong type")
	}

	privPemBytes := privPem.Bytes

	privateKey, err := x509.ParsePKCS1PrivateKey(privPemBytes)
	if err != nil {
		return nil, err
	}

	pubPem, _ := pem.Decode([]byte(t.PublicKey))

	if pubPem.Type != "PUBLIC KEY" {
		return nil, errors.New("public key is of the wrong type")
	}

	pubPemBytes := pubPem.Bytes

	parsedKey, err := x509.ParsePKIXPublicKey(pubPemBytes)
	if err != nil {
		return nil, err
	}

	pubKey, ok := parsedKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid public key")
	}

	privateKey.PublicKey = *pubKey

	t.rsaKey = privateKey
	return privateKey, nil
}

// UserRepoConfig contains options to connect to the user graphQL repo
type UserRepoConfig struct {
	// The user graphQL server URL
	Url string `json:"url"`
}

// Config all options required by this service to run
type Config struct {
	Token     Token          `json:"token"`
	UserRepo  UserRepoConfig `json:"userRepo"`
	Host      string         `json:"host,omitempty"`
	Port      string         `json:"port,omitempty"`
	APIPrefix string         `json:"apiPrefix,omitempty"`
}

// DefaultConfig returns a configuration object with the default values
func DefaultConfig() Config {
	return Config{
		Token: Token{
			Duration:        7 * 24 * 60 * 60,  // 7 days
			RefreshDuration: 30 * 24 * 60 * 60, // 30 days
		},
		Host:      "0.0.0.0",
		Port:      "8080",
		APIPrefix: "/auth",
	}
}

// LoadConfiguration reads configuration from the specified json file
func LoadConfiguration(file string) Config {
	config := DefaultConfig()
	configFile, err := os.Open(file)

	if err != nil {
		log.Warn().Err(err).Msg("Can't load config file. Default values will be used instead")
	}

	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
