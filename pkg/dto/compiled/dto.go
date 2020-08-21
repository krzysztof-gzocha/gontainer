package compiled

type Param struct {
	Code string
	Raw  interface{}
}

type Service struct {
	Name        string
	Getter      string
	Type        string
	Value       string
	Constructor string
	Args        []interface{}
}

type DTO struct {
	Meta struct {
		Pkg           string
		ContainerType string
	}
	Params   map[string]Param
	Services []Service
}
