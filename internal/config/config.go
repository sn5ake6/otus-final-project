package config

import (
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Logger     LoggerConf
	Storage    StorageConf
	HTTPServer HTTPServerConf
	GRPCServer GRPCServerConf
	Limit      LimitConf
}

func NewConfig() Config {
	return Config{}
}

func LoadConfig(configFile string, config interface{}) error {
	content, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(content, config)
	if err != nil {
		return err
	}

	return nil
}
