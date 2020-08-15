package input

import (
	"github.com/gomponents/gontainer/pkg/parameters"
)

type Service struct {
	Getter      string   `yaml:"getter"`
	Type        string   `yaml:"type"`
	Constructor string   `yaml:"constructor"`
	Disposable  bool     `yaml:"disposable"`
	Args        []string `yaml:"args"`
	Tags        []string `yaml:"tags"`
}

type DTO struct {
	Meta struct {
		Pkg           string            `yaml:"pkg"` // todo make default main
		ContainerType string            `yaml:"container_type"`
		Imports       map[string]string `yaml:"imports"`
		Functions     map[string]string `yaml:"functions"` // todo e.g. "env" => "os.Getenv"
	} `yaml:"meta"`
	Params   parameters.RawParameters `yaml:"parameters"`
	Services map[string]Service       `yaml:"services"`
}
