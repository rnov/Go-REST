package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func LoadAPIConfig(confPath string) (APIConfig, error) {
	var config APIConfig
	data, err := ioutil.ReadFile(confPath)
	if err != nil {
		log.Fatal("error reading configuration " + confPath)
	}

	err = yaml.Unmarshal(data, &config)

	if err != nil {
		return config, err
	}

	return config, nil
}

type DBConfig struct {
	Type string `yaml:"type"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	DB   int    `yaml:"dbInterface"`
}

type LoggerConfig struct {
	Name   string `yaml:"name"`
	File   string `yaml:"file"`
	Active bool   `yaml:"active"`
}

type APIConfig struct {
	Server   Server       `yaml:"server"`
	DBCfg    DBConfig     `yaml:"dbConfig"`
	LogsPath string       `yaml:"loggers_paths"`
	RedisLog LoggerConfig `yaml:"redis_logger"`
	//	... api, postgres, logger ...
}

type Server struct {
	Address string `yaml:"address"`
}
