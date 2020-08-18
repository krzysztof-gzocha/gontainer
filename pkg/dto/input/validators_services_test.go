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
			name:  "my..service",
			error: "service name must match pattern `" + regexServiceName.String() + "`, `my..service` given",
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

			if assert.Error(t, err) {
				assert.Equal(t, s.error, err.Error())
			}
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

			if assert.Error(t, err) {
				assert.Equal(t, s.error, err.Error())
			}
		})
	}
}

func TestValidateServiceGetter(t *testing.T) {
	scenarios := []struct {
		getter string
		error  string
	}{
		{
			getter: "",
		},
		{
			getter: "GetName",
		},
		{
			getter: "getName",
		},
		{
			getter: "0getName",
			error:  "getter must match `" + regexServiceGetter.String() + "`, `0getName` given",
		},
		{
			getter: "Get Name",
			error:  "getter must match `" + regexServiceGetter.String() + "`, `Get Name` given",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServiceGetter(Service{Getter: s.getter})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			if assert.Error(t, err) {
				assert.Equal(t, s.error, err.Error())
			}
		})
	}
}

func TestValidateServiceType(t *testing.T) {
	scenarios := []struct {
		type_ string
		error string
	}{
		{
			type_: "",
		},
		{
			type_: "my/import/foo.Bar",
		},
		{
			type_: "*my/import/foo.Bar",
		},
		{
			type_: "**my/import/foo.Bar",
			error: "type must match `" + regexServiceType.String() + "`, `**my/import/foo.Bar` given",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServiceType(Service{Type: s.type_})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			if assert.Error(t, err) {
				assert.Equal(t, s.error, err.Error())
			}
		})
	}
}

func TestValidateServiceValue(t *testing.T) {
	scenarios := []struct {
		val   string
		error string
	}{
		{
			val: "",
		},
		{
			val: "my/import/foo.Bar",
		},
		{
			val: "my/import/foo.MyStruct{}.Bar",
		},
		{
			val:   "my/import/foo.MyStruct{.Bar",
			error: "value must match `" + regexServiceValue.String() + "`, `my/import/foo.MyStruct{.Bar` given",
		},
		{
			val:   "my/import/foo",
			error: "value must match `" + regexServiceValue.String() + "`, `my/import/foo` given",
		},
		{
			val:   "*my/import/foo.Bar",
			error: "value must match `" + regexServiceValue.String() + "`, `*my/import/foo.Bar` given",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("Scenario #%d", i), func(t *testing.T) {
			err := ValidateServiceValue(Service{Value: s.val})
			if s.error == "" {
				assert.NoError(t, err)
				return
			}

			if assert.Error(t, err) {
				assert.Equal(t, s.error, err.Error())
			}
		})
	}
}
