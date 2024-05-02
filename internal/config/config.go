package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

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

type HTTPServer struct {
	Host                string        `yaml:"host"`
	Port                int           `yaml:"port"`
	ComixUpdateInterval time.Duration `yaml:"comix-update-interval"`
}

type Postgres struct {
	DSN        string `yaml:"dsn"`
	Migrations string `yaml:"migrations-path"`
}

type Config struct {
	JSONIndex    `yaml:"json-index"`
	JSONDatabase `yaml:"json-database"`
	XKCDCom      `yaml:"xkcd-com"`
	HTTPServer   `yaml:"http-server"`
	Postgres     `yaml:"postgres"`
}

var DefaultConfigPath string = func() string {
	currentFile, _ := filepath.Abs("./")
	root := findModuleRoot(currentFile)
	return path.Join(root, "config", "config.yaml")
}()

func findModuleRoot(dir string) string {
	// Source: https://github.com/golang/go/blob/9e3b1d53a012e98cfd02de2de8b1bd53522464d4/src/cmd/go/internal/modload/init.go#L1504
	dir = filepath.Clean(dir)

	for {
		if fi, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil && !fi.IsDir() {
			return dir
		}
		d := filepath.Dir(dir)
		if d == dir {
			break
		}
		dir = d
	}
	return ""
}

func (c *Config) makePathAbsolute() error {
	currentFile, err := filepath.Abs("./")
	if err != nil {
		return err
	}
	root := findModuleRoot(currentFile)

	c.JSONIndex.FileName = path.Join(root, c.JSONIndex.FileName)
	c.JSONDatabase.FileName = path.Join(root, c.JSONIndex.FileName)
	c.Postgres.Migrations = path.Join(root, c.Postgres.Migrations)
	return nil
}

func New(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}

	err = config.makePathAbsolute()
	return config, err
}
