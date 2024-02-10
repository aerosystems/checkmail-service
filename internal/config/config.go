package config

import "github.com/spf13/viper"

type Config struct {
	Mode         string `mapstructure:"MODE"`
	PostgresDSN  string `mapstructure:"POSTGRES_DSN"`
	AccessSecret string `mapstructure:"ACCESS_SECRET"`
}

func NewConfig() *Config {
	var cfg Config
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}
	return &cfg
}
