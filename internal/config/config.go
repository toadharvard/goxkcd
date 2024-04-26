package config

import (
	"fmt"
	"os"
	"time"

	"log/slog"

	"github.com/toadharvard/goxkcd/pkg/iso6391"
	"gopkg.in/yaml.v3"
)

type ISOCode6391 iso6391.ISOCode6391

func (code *ISOCode6391) UnmarshalYAML(value *yaml.Node) error {
	var str string
	if err := value.Decode(&str); err != nil {
		return fmt.Errorf("cannot decode ISO 639-1 code: %w", err)
	}

	parsed, err := iso6391.NewLanguage(str)
	if err != nil {
		return err
	}

	*code = ISOCode6391(parsed)
	return nil
}

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
	Language        ISOCode6391   `yaml:"language"`
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

	slog.Debug("config loaded", "config", config)
	return config, nil
}
