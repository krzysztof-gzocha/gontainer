package input

type Validator interface {
	Validate(DTO) error
}

type ChainValidator struct {
	validators []func(DTO) error
}

func NewChainValidator(validators ...func(DTO) error) *ChainValidator {
	return &ChainValidator{validators: validators}
}

func (c ChainValidator) Validate(m DTO) error {
	for _, v := range c.validators {
		if err := v(m); err != nil {
			return err
		}
	}
	return nil
}

func NewDefaultValidator() Validator {
	validators := make([]func(DTO) error, 0)
	validators = append(validators, DefaultMetaValidators()...)
	return NewChainValidator(validators...)
}
