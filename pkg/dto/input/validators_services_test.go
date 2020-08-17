package input

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateServiceName(t *testing.T) {
	scenarios := []struct {
		name  string
		error string
	}{
		{
			name: "my.service",
		},
		{
			name:  "%service%",
			error: "service name must match pattern `" + regexServiceName.String() + "`, `%service%` given",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServiceName(s.name)
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.Equal(t, s.error, err.Error())
		})
	}
}

func TestValidateConstructorType(t *testing.T) {
	scenarios := []struct {
		service Service
		error   string
	}{
		{
			service: Service{
				Constructor: "NewService",
			},
		},
		{
			error: "missing constructor, value or type",
		},
		{
			service: Service{
				Value:       "MyValue",
				Constructor: "MyConstructor",
			},
			error: "cannot define constructor and value together",
		},
		{
			service: Service{
				Constructor: "MyConstructor",
				Type:        "MyType",
			},
			error: "defined type will not be used, provide getter",
		},
		{
			service: Service{
				Value: "MyValue",
				Type:  "MyType",
			},
			error: "defined type will not be used, provide getter",
		},
		{
			service: Service{
				Type: "MyType",
				Args: []string{"param"},
			},
			error: "arguments are not empty, but constructor is missing",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateConstructorType(s.service)
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			assert.Error(t, err)
			assert.Equal(t, s.error, err.Error())
		})
	}
}
