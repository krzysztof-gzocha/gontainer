package compiled

type Param struct {
	Code string
	Raw  interface{}
}

type DTO struct {
	Meta struct {
		Pkg           string
		ContainerType string
	}
	Params map[string]Param
}
