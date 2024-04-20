package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type ISOCode639_1 = string

type JSONDatabase struct {
	FileName string `yaml:"file-name"`
}

type JSONIndex struct {
	FileName string `yaml:"file-name"`
}

type XKCDCom struct {
	URL             string        `yaml:"url"`
	BatchSize       int           `yaml:"batch-size"`
	NumberOfWorkers int           `yaml:"number-of-workers"`
	Language        ISOCode639_1  `yaml:"language"`
	Timeout         time.Duration `yaml:"timeout"`
}

type Config struct {
	JSONIndex    `yaml:"json-index"`
	JSONDatabase `yaml:"json-database"`
	XKCDCom      `yaml:"xksd-com"`
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
