package userconfig

import (
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	MiseGoPath      string `mapstructure:"MISE_GO_PATH"`
	DefaultFilePerm string `mapstructure:"DEFAULT_FILE_PERMISSION"`
}

func LoadConfig() (AppConfig, error) {
	var config AppConfig

	viper.AutomaticEnv()

	if err := viper.BindEnv("MISE_GO_PATH"); err != nil {
		return AppConfig{}, fmt.Errorf("error reading config %v", err)
	}

	if err := viper.BindEnv("DEFAULT_FILE_PERMISSION"); err != nil {
		return AppConfig{}, fmt.Errorf("error reading config %v", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, fmt.Errorf("failed to unmarshal env config: %w", err)
	}

	return config, nil
}
