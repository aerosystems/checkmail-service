package config

import (
	"github.com/spf13/viper"
)

const (
	defaultMode    = "prod"
	defaultWebPort = 8080
)

type Config struct {
	Mode                         string
	WebPort                      int
	GcpProjectId                 string
	GoogleApplicationCredentials string
	PostgresDSN                  string
	ProjectServiceRpcAddress     string
	AccessSecret                 string
}

func NewConfig() *Config {
	viper.AutomaticEnv()
	mode := viper.GetString("CHCKML_MODE")
	if mode == "" {
		mode = defaultMode
	}
	webPort := viper.GetInt("PORT")
	if webPort == 0 {
		webPort = defaultWebPort
	}
	return &Config{
		Mode:                         mode,
		WebPort:                      webPort,
		GcpProjectId:                 viper.GetString("GCP_PROJECT_ID"),
		GoogleApplicationCredentials: viper.GetString("GOOGLE_APPLICATION_CREDENTIALS"),
		PostgresDSN:                  viper.GetString("POSTGRES_DSN"),
		ProjectServiceRpcAddress:     viper.GetString("PROJECT_SERVICE_RPC_ADDR"),
		AccessSecret:                 viper.GetString("ACCESS_SECRET"),
	}
}
