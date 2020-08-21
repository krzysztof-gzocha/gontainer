package compiled

type Validator struct {
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v Validator) Validate(DTO) error {
	// todo
	return nil
}

func NewDefaultValidator() *Validator {
	// todo
	return NewValidator()
}
