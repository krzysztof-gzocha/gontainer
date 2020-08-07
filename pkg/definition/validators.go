package definition

type Validator interface {
	Validate(Definition) error
}

type ChainValidator struct {
	validators []func(Definition) error
}

func NewChainValidator(validators []func(Definition) error) *ChainValidator {
	return &ChainValidator{validators: validators}
}

func (c ChainValidator) Validate(d Definition) error {
	for _, v := range c.validators {
		err := v(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewDefaultValidator() Validator {
	return NewChainValidator([]func(Definition) error{
		ValidateMetaPkg,
		ValidateMetaImports,
		ValidateMetaContainerType,
		ValidateParams,
		ValidateServicesNames,
	})
}
