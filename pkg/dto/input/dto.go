package input

import (
	"fmt"
)

const (
	defaultPkg           = "main"
	defaultContainerType = "Gontainer"
)

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

type Tag struct {
	Name     string
	Priority int
}

func (t *Tag) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var z interface{}
	if err := unmarshal(&z); err != nil {
		return nil
	}

	if s, ok := z.(string); ok {
		t.Name = s
		t.Priority = 0
		return nil
	}

	if m, ok := z.(map[interface{}]interface{}); ok {
		if n, ok := m["name"]; ok {
			name, okName := n.(string)
			if !okName {
				return fmt.Errorf("name must be an instance of string")
			}
			t.Name = name
		} else {
			return fmt.Errorf("name of tag is required")
		}

		if p, ok := m["priority"]; ok {
			priority, okPriority := p.(int)
			if !okPriority {
				return fmt.Errorf("priority must be an instance of int")
			}
			t.Priority = priority
		} else {
			t.Priority = 0
		}
	}
	return nil
}

type Service struct {
	Getter      string                 `yaml:"getter"`      // e.g. GetDB
	Type        string                 `yaml:"type"`        // *?my/import/path.Type
	Value       string                 `yaml:"value"`       // my/import/path.Variable
	Constructor string                 `yaml:"constructor"` // NewDB
	Args        []interface{}          `yaml:"args"`        // ["%host%", "%port%", "@logger"]
	Calls       []Call                 `yaml:"calls"`       // [["SetLogger", ["@logger"]], ...]
	Fields      map[string]interface{} `yaml:"fields"`      // Field: "%value%"
	Tags        []Tag                  `yaml:"tags"`        // ["service_decorator", ...]
	Disposable  bool                   `yaml:"disposable"`  // if true container creates new instance of given service always, otherwise service is cached
	Todo        bool                   `yaml:"todo"`        // if true skips validation and returns error whenever users asks container for a service
}

type DTO struct {
	Meta struct {
		Pkg           string            `yaml:"pkg"`            // default "main"
		ContainerType string            `yaml:"container_type"` // default "Gontainer"
		Imports       map[string]string `yaml:"imports"`        // [["alias": "my/long/path"], ...]
		Functions     map[string]string `yaml:"functions"`      // [["env": "os.Getenv"], ...]
	} `yaml:"meta"`
	Params   map[string]interface{} `yaml:"parameters"`
	Services map[string]Service     `yaml:"services"`
}

func CreateDefaultDTO() DTO {
	result := DTO{}
	result.Meta.Pkg = defaultPkg
	result.Meta.ContainerType = defaultContainerType
	result.Meta.Functions = map[string]string{
		"env":    "os.Env",
		"envInt": "github.com/gomponents/gontainer-helpers/env.MustGetInt",
		"todo":   "github.com/gomponents/gontainer-helpers/std.GetMissingParameter",
	}
	return result
}
