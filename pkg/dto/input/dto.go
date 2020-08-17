package input

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg/parameters"
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
	fmt.Println(z)

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
	Getter      string   `yaml:"getter"`
	Type        string   `yaml:"type"`
	Value       string   `yaml:"value"`
	Constructor string   `yaml:"constructor"`
	Args        []string `yaml:"args"`
	Calls       []Call   `yaml:"calls"`
	Disposable  bool     `yaml:"disposable"`
	Tags        []string `yaml:"tags"`
	Todo        bool     `yaml:"todo"`
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
