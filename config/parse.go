package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/goplus/llgo/xtool/env"
)

const LLCPPG_CFG = "llcppg.cfg"
const LLCPPG_SYMB = "llcppg.symb.json"
const LLCPPG_SIGFETCH = "llcppg.sigfetch.json"
const LLCPPG_PUB = "llcppg.pub"

// json middleware for validating
func (c *Config) UnmarshalJSON(data []byte) error {
	// create a new type here to avoid unmarshalling infinite loop.
	type newConfig Config

	var config newConfig
	err := json.Unmarshal(data, &config)

	if err != nil {
		return err
	}

	*c = Config(config)

	// do some check

	// when headeronly mode is disabled, libs must not be empty.
	if c.Libs == "" && !c.HeaderOnly {
		return fmt.Errorf("%w: libs must not be empty", ErrConfig)
	}

	return nil
}

func ParseConfigFile(cfgFile string) (*Config, error) {
	openCfgFile, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}
	defer openCfgFile.Close()
	var conf Config
	err = json.NewDecoder(openCfgFile).Decode(&conf)
	if err != nil {
		return nil, err
	}
	conf.CFlags = env.ExpandEnv(conf.CFlags)
	conf.Libs = env.ExpandEnv(conf.Libs)
	return &conf, nil
}

func MarshalConfigFile(cfgFile string) ([]byte, error) {
	conf, err := ParseConfigFile(cfgFile)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(&conf, "", "  ")
}
