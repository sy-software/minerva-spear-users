package domain

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
)

type Config struct {
}

func DefaultConfig() Config {
	return Config{}
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
