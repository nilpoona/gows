package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server ServerConfig
}

type ServerConfig struct {
	Host string
	Port string
}

func NewConfig() *Config {
	var config *Config
	env := os.Getenv("GO_ENV")
	_, err := toml.DecodeFile("conf/"+env+".toml", &config)
	if err != nil {
		log.Printf("error: %v", err)
	}

	return config
}
