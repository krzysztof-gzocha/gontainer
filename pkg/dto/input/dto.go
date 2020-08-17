package input

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg/parameters"
)

const (
	defaultPkg           = "main"
	defaultContainerType = "Gontainer"
)

// todo https://symfony.com/blog/new-in-symfony-4-3-configuring-services-with-immutable-setters
type Call struct {
	Method    string
	Args      []interface{}
	Immutable bool
}

func (c *Call) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var z []interface{}
	if err := unmarshal(&z); err != nil {
		return err
	}

	if len(z) == 0 || len(z) > 3 {
		return fmt.Errorf("object Call must contain 1 - 3 args, %d given", len(z))
	}

	if s, ok := z[0].(string); !ok {
		return fmt.Errorf("first element of object Call must be a string, `%T` given", z[0])
	} else {
		c.Method = s
	}

	if len(z) >= 2 {
		if args, ok := z[1].([]interface{}); !ok {
			return fmt.Errorf("second element of object Call must be an array, `%T` given", z[1])
		} else {
			c.Args = args
		}
	}

	if len(z) >= 3 {
		if i, ok := z[2].(bool); !ok {
			return fmt.Errorf("third element of object Call must be a bool, `%T` given", z[2])
		} else {
			c.Immutable = i
		}
	}

	return nil
}

type Service struct {
	Getter      string                 `yaml:"getter"`      // e.g. GetDB
	Type        string                 `yaml:"type"`        // *?my/import/path.Type
	Value       string                 `yaml:"value"`       // my/import/path.Variable
	Constructor string                 `yaml:"constructor"` // NewDB
	Args        []string               `yaml:"args"`        // ["%host%", "%port%", "@logger"]
	Calls       []Call                 `yaml:"calls"`       // [["SetLogger", ["@logger"]], ...]
	Fields      map[string]interface{} `yaml:"fields"`      // Field: "%value%"
	Disposable  bool                   `yaml:"disposable"`  // if true container creates new instance of given service always, otherwise service is cached
	Tags        []string               `yaml:"tags"`        // ["service_decorator", ...]
	Todo        bool                   `yaml:"todo"`        // if true skips validation and returns error whenever users asks container for a service
}

type DTO struct {
	Meta struct {
		Pkg           string            `yaml:"pkg"`            // default "main"
		ContainerType string            `yaml:"container_type"` // default "Gontainer"
		Imports       map[string]string `yaml:"imports"`
		Functions     map[string]string `yaml:"functions"`
	} `yaml:"meta"`
	Params   parameters.RawParameters `yaml:"parameters"`
	Services map[string]Service       `yaml:"services"`
}

func CreateDefaultDTO() DTO {
	result := DTO{}
	result.Meta.Pkg = defaultPkg
	result.Meta.ContainerType = defaultContainerType
	result.Meta.Functions = map[string]string{
		"os.Env": "env",
		"envInt": "github.com/gomponents/gontainer-helpers/env.MustGetInt",
		"todo":   "github.com/gomponents/gontainer-helpers/std.GetMissingParameter",
	}
	return result
}
