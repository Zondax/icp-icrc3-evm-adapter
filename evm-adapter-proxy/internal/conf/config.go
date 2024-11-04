package conf

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Logging      LoggingConfig `mapstructure:"logging"`
	Metrics      MetricsConfig `mapstructure:"metrics"`
	RouterConfig RouterConfig  `mapstructure:"routerConfig"`
	ServerPort   string        `mapstructure:"serverPort"`
	ICP          *ICPConfig    `mapstructure:"icp"`
}

type LoggingConfig struct {
	Level string `mapstructure:"level"`
}

type MetricsConfig struct {
	Path                  string `mapstructure:"path"`
	Port                  string `mapstructure:"port"`
	SystemMetricsInterval string `mapstructure:"systemMetricsInterval"`
}

type RouterConfig struct {
	AppRevision string
	AppVersion  string
}

type ICPConfig struct {
	LoggerCanisterID string `mapstructure:"loggerCanisterId"`
	DexCanisterID    string `mapstructure:"dexCanisterId"`
	NodeURL          string `mapstructure:"nodeUrl"`
}

func (c Config) SetDefaults() {
	viper.SetDefault("icp.canisterId", "")
	viper.SetDefault("icp.nodeUrl", "https://ic0.app")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("metrics.path", "/metrics")
	viper.SetDefault("metrics.port", "9090")
}

func (c Config) Validate() error {
	if c.ICP.LoggerCanisterID == "" {
		return fmt.Errorf("ICP LoggerCanisterID must be provided")
	}
	if c.ICP.DexCanisterID == "" {
		return fmt.Errorf("ICP DexCanisterID must be provided")
	}
	if c.ICP.NodeURL == "" {
		return fmt.Errorf("ICP NodeURL must be provided")
	}
	return nil
}
