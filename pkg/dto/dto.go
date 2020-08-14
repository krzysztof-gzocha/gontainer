package dto

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

type Input struct {
	Meta struct {
		Pkg           string            `yaml:"pkg"`
		Imports       map[string]string `yaml:"imports"`
		ContainerType string            `yaml:"container_type"`
	} `yaml:"meta"`
	Params   parameters.RawParameters `yaml:"parameters"`
	Services map[string]Service       `yaml:"services"`
}

type ServiceLink struct {
	Name string
	Type string
}

type CompiledArg struct {
	Code        string
	Raw         string
	ServiceLink *ServiceLink
	DependsOn   []string // list of depending parameters todo validation
}

func (c *CompiledArg) IsService() bool {
	return c.ServiceLink != nil
}

type CompiledService struct {
	Name        string
	Getter      string
	Value       string // todo Value cannot be use with getter
	Type        string
	Constructor string
	Disposable  bool
	Args        []CompiledArg
	Tags        []string
}

// TODO given model will be passed to template
type CompiledInput struct {
	Meta struct {
		Pkg           string
		ContainerType string
	}
	Params   parameters.ResolvedParams
	Services []CompiledService
}
