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
	// Token duration in seconds
	Duration int64
	// Refresh Token in seconds
	RefreshDuration int64
	// For JWT signature using RS256 algorithm
	PrivateKey string
	PublicKey  string
}

func (t *Token) KeyPair() (*rsa.PrivateKey, error) {
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

	return privateKey, nil
}

type Config struct {
	Token Token
}

func DefaultConfig() Config {
	return Config{
		Token: Token{
			Duration:        7 * 24 * 60 * 60,  // 7 days
			RefreshDuration: 30 * 24 * 60 * 60, // 30 days
		},
	}
}

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
