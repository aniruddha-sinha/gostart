package userconfig

import (
	"bytes"
	"embed"
	"fmt"

	"github.com/spf13/viper"
)

//go:embed defaults.yaml
var defaultConfigFile embed.FS

type AppConfig struct {
	MiseGoPath      string `mapstructure:"mise_go_path"`
	DefaultFilePerm string `mapstructure:"default_file_permission"`
}

func LoadConfig() (AppConfig, error) {
	var config AppConfig

	defaultBytes, err := defaultConfigFile.ReadFile("defaults.yaml")
	if err != nil {
		return config, fmt.Errorf("failed to read embedded defaults: %w", err)
	}

	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewReader(defaultBytes)); err != nil {
		return config, fmt.Errorf("failed to parse embedded defaults: %w", err)
	}

	viper.AutomaticEnv()
	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
