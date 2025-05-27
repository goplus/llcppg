package config

import (
	"encoding/json"
	"io"
	"os"
)

func GetConfFromStdin() (conf Config, err error) {
	return ConfigFromReader(os.Stdin)
}

func GetConfFromFile(cfgFile string) (conf Config, err error) {
	fileReader, err := os.Open(cfgFile)
	if err != nil {
		return
	}
	defer fileReader.Close()

	return ConfigFromReader(fileReader)
}

func ConfigFromReader(reader io.Reader) (Config, error) {
	var config Config

	if err := json.NewDecoder(reader).Decode(&config); err != nil {
		return Config{}, err
	}

	return config, nil
}
