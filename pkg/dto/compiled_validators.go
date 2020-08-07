package dto

type CompiledValidator interface {
	Validate(CompiledInput) error
}

type ChainCompiledValidator struct {
	validators []func(CompiledInput) error
}

func NewChainCompiledValidator(validators []func(CompiledInput) error) *ChainCompiledValidator {
	return &ChainCompiledValidator{validators: validators}
}

func (c *ChainCompiledValidator) Validate(i CompiledInput) error {
	for _, v := range c.validators {
		if err := v(i); err != nil {
			return err
		}
	}
	return nil
}

func NewDefaultCompiledValidator() CompiledValidator {
	return NewChainCompiledValidator([]func(CompiledInput) error{
		validateCircularDependency,
	})
}

func validateCircularDependency(i CompiledInput) error {
	deps := make(map[string][]string)
	for n, s := range i.Services {
		deps[n] = make([]string, 0)
		for _, a := range s.Args {
			if !a.IsService() {
				continue
			}
			deps[n] = append(deps[n], a.ServiceLink.Name)
		}
	}

	// todo

	return nil
}
