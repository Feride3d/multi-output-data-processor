package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Http string
	Conf []struct {
		Tag     string        `yaml:"tag"`
		Outputs []interface{} `yaml:"outputs"`
	} `yaml:"conf"`
}

func NewConfig(file string) (conf Config, err error) {
	viper.AutomaticEnv()
	if file != "" {
		viper.SetConfigFile(file)
		viper.SetConfigType("yml")
		if err := viper.ReadInConfig(); err != nil {
			return Config{}, fmt.Errorf("reading config error: %w", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return Config{}, fmt.Errorf("unmarshaling config error: %w", err)
	}

	return cfg, nil
}
