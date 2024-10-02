package conf

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"github.com/zondax/golem/pkg/zrouter"
)

type Config struct {
	Logging      LoggingConfig  `yaml:"logging"`
	Metrics      MetricsConfig  `yaml:"metrics"`
	RouterConfig zrouter.Config `yaml:"routerConfig"`
	ServerPort   string         `yaml:"serverPort"`
	ICP          *ICPConfig     `yaml:"icp"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
}

type MetricsConfig struct {
	Path                  string        `yaml:"path"`
	Port                  string        `yaml:"port"`
	SystemMetricsInterval time.Duration `yaml:"systemMetricsInterval"`
}

type ICPConfig struct {
	CanisterID string `yaml:"canisterId"`
	NodeURL    string `yaml:"nodeUrl"`
}

func (c Config) SetDefaults() {
	viper.SetDefault("icp.canisterId", "")
	viper.SetDefault("icp.nodeUrl", "https://ic0.app")
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("metrics.path", "/metrics")
	viper.SetDefault("metrics.port", "9090")
}

func (c Config) Validate() error {
	if c.ICP.CanisterID == "" {
		return fmt.Errorf("ICP CanisterID must be provided")
	}
	if c.ICP.NodeURL == "" {
		return fmt.Errorf("ICP NodeURL must be provided")
	}
	return nil
}
