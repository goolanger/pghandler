package src

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Server struct {
		Host         string `yaml:"host", envconfig:"SERVER_HOST"`
		Port         string `yaml:"port", envconfig:"SERVER_PORT"`
		ReadTimeout  string `yaml:"readtimeout"`
		WriteTimeout string `yaml:"writetimeout"`
	} `yaml:"server"`
	Instance []struct {
		Name     string `yaml:"name"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Database []struct {
			Name     string `yaml:"name"`
			Schema   string `yaml:"schema"`
			Username string `yaml:"user"`
			Password string `yaml:"pass"`
		} `yaml:"databases"`
	} `yaml:"instances"`
}

func ReadConfiguration() (*Config, error) {
	f, err := os.Open("config/application.yml")
	if err != nil {
		processError(err)
	}
	defer f.Close()

	var cfg Config
	cfg.Server.ReadTimeout = "10s"
	cfg.Server.WriteTimeout = "10s"
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	return &cfg, err
}

func processError(err error) {
	log.Fatal(err)
}
