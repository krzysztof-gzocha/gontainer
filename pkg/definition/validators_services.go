package definition

import (
	"fmt"
	"regexp"

	"github.com/gomponents/gontainer/pkg/syntax"
)

var (
	serviceNameRegex   = regexp.MustCompile(`^` + syntax.ServiceNamePattern + `$`)
	serviceGetterRegex = regexp.MustCompile(`^[A-Z][A-Za-z0-9_]$`)
)

func ValidateServicesNames(d Definition) error {
	for n, _ := range d.Services {
		if !serviceNameRegex.MatchString(n) {
			return fmt.Errorf(
				"service name must match pattern `%s`, `%s` given",
				serviceNameRegex.String(),
				n,
			)
		}
	}

	return nil
}

func ValidateServicesGetters(d Definition) error {
	for _, s := range d.Services {
		if s.Getter != "" && !serviceGetterRegex.MatchString(s.Getter) {
			return fmt.Errorf(
				"getter must match `%s`, `%s` given",
				serviceGetterRegex.String(),
				s.Getter,
			)
		}
	}

	return nil
}
