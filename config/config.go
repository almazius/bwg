package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Server struct {
		Port    string `json:"port"`
		Timeout int    `json:"timeout"`
	} `json:"server"`
	Postgres struct {
		User     string `json:"user"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     string `json:"port"`
		DbName   string `json:"dbName"`
	}
}

func LoadConfig() (*viper.Viper, error) {
	v := viper.New()
	v.AddConfigPath("config")
	v.SetConfigName("config")
	v.SetConfigType("yml")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("unable to read config, %v", err)
		return nil, err
	}
	return v, nil
}

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config
	err := v.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode config into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
