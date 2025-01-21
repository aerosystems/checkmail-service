package config

import (
	"github.com/spf13/viper"
)

const (
	defaultMode  = "prod"
	defaultPort  = 8080
	defaultProto = "http"
)

type Config struct {
	Mode                         string
	Port                         int
	Proto                        string
	GcpProjectId                 string
	GoogleApplicationCredentials string
	PostgresDSN                  string
}

func NewConfig() *Config {
	viper.AutomaticEnv()
	mode := viper.GetString("MODE")
	if mode == "" {
		mode = defaultMode
	}
	port := viper.GetInt("PORT")
	if port == 0 {
		port = defaultPort
	}
	proto := viper.GetString("PROTO")
	if proto == "" {
		proto = defaultProto
	}
	return &Config{
		Mode:                         mode,
		Port:                         port,
		Proto:                        proto,
		GcpProjectId:                 viper.GetString("GCP_PROJECT_ID"),
		GoogleApplicationCredentials: viper.GetString("GOOGLE_APPLICATION_CREDENTIALS"),
		PostgresDSN:                  viper.GetString("POSTGRES_DSN"),
	}
}
