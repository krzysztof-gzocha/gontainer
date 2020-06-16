package pkg

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBaseParamContainer(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]interface{}
		hasNot []string
	}{
		{
			name:   "empty container",
			params: nil,
			hasNot: []string{"param"},
		},
		{
			name:   "with parameters",
			params: map[string]interface{}{"one": 1, "pi": math.Pi, "name": "gontainer"},
			hasNot: []string{"e"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewBaseParamContainer(tt.params)
			for id, val := range tt.params {
				assert.True(t, container.HasParam(id))
				assert.Equal(t, val, container.MustGetParam(id))

				realValue, err := container.GetParam(id)
				assert.NoError(t, err)
				assert.Equal(t, val, realValue)
			}
			for _, id := range tt.hasNot {
				assert.False(t, container.HasParam(id))
				assert.Panics(t, func() {
					container.MustGetParam(id)
				})

				realValue, err := container.GetParam(id)
				assert.Error(t, err)
				assert.Nil(t, realValue)
			}
		})
	}
}

func TestNewBaseContainer(t *testing.T) {
	tests := []struct {
		name    string
		getters map[string]GetterDefinition
		hasNot  []string
	}{
		{
			name:    "empty container",
			getters: nil,
			hasNot:  []string{"db"},
		},
		{
			name: "db",
			getters: map[string]GetterDefinition{
				"db": createGetterDefinition(struct{}{}, "", false),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewBaseContainer(tt.getters)

			for id, _ := range tt.getters {
				assert.True(t, container.Has(id))
				assert.NotPanics(t, func() {
					container.MustGet(id)
				})

				_, err := container.Get(id)
				assert.NoError(t, err)
			}

			for _, id := range tt.hasNot {
				assert.False(t, container.Has(id))
				assert.Panics(t, func() {
					container.MustGet(id)
				})

				_, err := container.Get(id)
				assert.Error(t, err)
			}
		})
	}
}

func createGetterDefinition(svc interface{}, err string, disposable bool) GetterDefinition {
	return GetterDefinition{
		Getter: func() (interface{}, error) {
			var e error
			if err != "" {
				e = fmt.Errorf(err)
			}

			return svc, e
		},
		Disposable: disposable,
	}
}

func TestNewBaseTaggedContainer(t *testing.T) {
	container := NewBaseContainer(map[string]GetterDefinition{
		"mysql": createGetterDefinition(struct{}{}, "", false),
		"redis": createGetterDefinition(struct{}{}, "", false),
	})

	tests := []struct {
		name      string
		container Container
		mapping   map[string][]string
		tags      map[string]uint
	}{
		{
			name:      "simple example",
			container: container,
			mapping:   map[string][]string{"storage": {"mysql", "redis"}},
			tags:      nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			container := NewBaseTaggedContainer(tt.container, tt.mapping)
			for tag, count := range tt.tags {
				services, err := container.GetByTag(tag)
				assert.NoError(t, err)
				assert.Equal(t, count, len(services))
			}
		})
	}
}
