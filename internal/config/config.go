package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ISOCode639_1 = string

type JSONDatabase struct {
	FileName string `yaml:"file-name"`
}

type XkcdCom struct {
	URL      string       `yaml:"url"`
	Language ISOCode639_1 `yaml:"language"`
}

type Config struct {
	JSONDatabase `yaml:"json-database"`
	XkcdCom      `yaml:"xksd-com"`
}

func New(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
