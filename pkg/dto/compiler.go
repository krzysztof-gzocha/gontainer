package dto

type Compiler interface {
	Compile(Input) CompiledInput
}

type BaseCompiler struct {
}

func (BaseCompiler) Compile(i Input) CompiledInput {
	r := CompiledInput{}
	r.Meta.Pkg = i.Meta.Pkg
	r.Meta.ContainerType = i.Meta.ContainerType
	// todo

	return r
}
