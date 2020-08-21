package pkg

import (
	"fmt"
	"io/ioutil"

	"github.com/gomponents/gontainer/pkg/dto/input"
	"gopkg.in/yaml.v2"
)

type ConfigReader interface {
	Read([]string) (input.DTO, error)
}

type simpleConfigReader struct {
	finder          finder
	beforeParseFile func(string)
}

func NewDefaultConfigReader(beforeParseFile func(string)) ConfigReader {
	return &simpleConfigReader{
		finder:          simpleFinder{},
		beforeParseFile: beforeParseFile,
	}
}

func (s simpleConfigReader) Read(patterns []string) (input.DTO, error) {
	files, err := s.finder.find(patterns)
	if err != nil {
		return input.DTO{}, err
	}

	if len(files) == 0 {
		return input.DTO{}, fmt.Errorf("cannot find any configuration file")
	}

	result := input.DTO{}

	for _, f := range files {
		s.beforeParseFile(f)
		if file, err := ioutil.ReadFile(f); err != nil {
			return input.DTO{}, fmt.Errorf("error has occurred during opening file `%s`: %s", f, err.Error())
		} else {
			if yamlErr := yaml.Unmarshal(file, &result); yamlErr != nil {
				return input.DTO{}, fmt.Errorf("error has occurred during parsing yaml file `%s`: %s", f, yamlErr.Error())
			}
		}
	}

	return result, nil
}
