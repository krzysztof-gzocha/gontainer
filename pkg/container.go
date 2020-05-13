package pkg

import (
	"fmt"
)

type Container interface {
	Get(string) (interface{}, error)
	MustGet(string) interface{}
	Has(string) bool
}

type TaggedContainer interface {
	Container
	GetByTag(string) ([]interface{}, error)
	MustGetByTag(string) []interface{}
}

type ParamContainer interface {
	GetParam(string) (interface{}, error)
	MustGetParam(string) interface{}
	HasParam(string) bool
}

type BaseParamContainer struct {
	params map[string]interface{}
}

func NewBaseParamContainer(params map[string]interface{}) *BaseParamContainer {
	return &BaseParamContainer{params: params}
}

func (b BaseParamContainer) GetParam(id string) (interface{}, error) {
	if b.HasParam(id) {
		return b.MustGetParam(id), nil
	}

	return nil, fmt.Errorf("parameter %s does not exist", id)
}

func (b BaseParamContainer) MustGetParam(id string) interface{} {
	return b.params[id]
}

func (b BaseParamContainer) HasParam(id string) bool {
	_, ok := b.params[id]
	return ok
}

type Getter func() (interface{}, error)

type GetterDefinition struct {
	Getter     Getter
	Disposable bool
}

type metaGetterDefintion struct {
	getter     Getter
	service    interface{}
	created    bool
	disposable bool
}

type BaseContainer struct {
	getters map[string]metaGetterDefintion
}

func NewBaseContainer(getters map[string]GetterDefinition) *BaseContainer {
	meta := make(map[string]metaGetterDefintion)
	for n, v := range getters {
		meta[n] = metaGetterDefintion{
			getter:     v.Getter,
			service:    nil,
			created:    false,
			disposable: v.Disposable,
		}
	}

	return &BaseContainer{getters: meta}
}

func (b BaseContainer) Get(id string) (interface{}, error) {
	if !b.Has(id) {
		return nil, fmt.Errorf("service %s does not exist", id)
	}

	getter := b.getters[id]
	if getter.created {
		return getter.service, nil
	}

	callee := getter.getter
	service, err := callee()

	if err != nil {
		return nil, fmt.Errorf("cannot create service %s: %s", id, err.Error())
	}

	if !getter.disposable {
		getter.created = true
		getter.service = service
		b.getters[id] = getter
	}

	return service, nil
}

func (b BaseContainer) MustGet(id string) interface{} {
	r, e := b.Get(id)

	if e != nil {
		panic(e)
	}

	return r
}

func (b BaseContainer) Has(id string) bool {
	_, ok := b.getters[id]
	return ok
}

type BaseTaggedContainer struct {
	container Container
	mapping map[string][]string
}

func (b BaseTaggedContainer) Get(id string) (interface{}, error) {
	return b.container.Get(id)
}

func (b BaseTaggedContainer) MustGet(id string) interface{} {
	return b.container.MustGet(id)
}

func (b BaseTaggedContainer) Has(id string) bool {
	return b.container.Has(id)
}

func (b BaseTaggedContainer) GetByTag(tag string) ([]interface{}, error) {
	result := make([]interface{}, 0)
	for _, id := range b.mapping[tag] {
		s, e := b.container.Get(id)
		if e != nil {
			return nil, fmt.Errorf("cannot get services by tag %s due to: %s", tag, e.Error())
		}
		result = append(result, s)
	}
	return result, nil
}

func (b BaseTaggedContainer) MustGetByTag(tag string) []interface{} {
	r, e := b.GetByTag(tag)
	if e != nil {
		panic(e)
	}
	return r
}

func NewBaseTaggedContainer(container Container, mapping map[string][]string) *BaseTaggedContainer {
	return &BaseTaggedContainer{container: container, mapping: mapping}
}
