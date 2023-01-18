package main

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Config struct {
	Database DBConfig
}

func ReadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	var conf Config
	err = toml.NewDecoder(file).Decode(&conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}
