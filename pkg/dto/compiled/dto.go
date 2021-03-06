package compiled

type Param struct {
	Name      string
	Code      string
	Raw       interface{}
	DependsOn []string
}

type Arg struct {
	Code              string
	Raw               interface{}
	DependsOnParams   []string
	DependsOnServices []string
	DependsOnTags     []string
}

type Call struct {
	Method    string
	Args      []Arg
	Immutable bool
}

type Field struct {
	Name  string
	Value Arg
}

type Tag struct {
	Name     string
	Priority int
}

type Service struct {
	Name        string
	Getter      string
	Type        string
	Value       string
	Constructor string
	Args        []Arg
	Calls       []Call
	Fields      []Field
	Tags        []Tag
	Disposable  bool
	Todo        bool
}

type DTO struct {
	Meta struct {
		Pkg           string
		ContainerType string
	}
	Params   []Param
	Services []Service
}
