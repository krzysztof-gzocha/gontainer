package definition

import (
	"github.com/gomponents/gontainer/pkg/parameters"
)

type Service struct {
	Getter      string   `yaml:"getter"`
	Type        string   `yaml:"type"`
	Constructor string   `yaml:"constructor"`
	WithError   bool     `yaml:"with_error"`
	Disposable  bool     `yaml:"disposable"`
	Args        []string `yaml:"args"`
	Tags        []string `yaml:"tags"`
}

type Definition struct {
	Meta struct {
		Pkg           string            `yaml:"pkg"`
		Imports       map[string]string `yaml:"imports"`
		ContainerType string            `yaml:"container_type"`
	} `yaml:"meta"`
	Params   parameters.RawParameters `yaml:"parameters"`
	Services map[string]Service       `yaml:"services"`
}
