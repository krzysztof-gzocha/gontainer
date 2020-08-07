package dto

type Validator interface {
	Validate(Input) error
}

type ChainValidator struct {
	validators []func(Input) error
}

func NewChainValidator(validators []func(Input) error) *ChainValidator {
	return &ChainValidator{validators: validators}
}

func (c *ChainValidator) Validate(i Input) error {
	for _, v := range c.validators {
		if err := v(i); err != nil {
			return err
		}
	}

	return nil
}

func NewDefaultValidator() Validator {
	return NewChainValidator([]func(Input) error{
		ValidateMetaPkg,
		ValidateMetaImports,
		ValidateMetaContainerType,
		ValidateParams,
		ValidateServices,
	})
}
