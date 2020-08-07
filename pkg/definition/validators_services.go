package definition

import (
	"fmt"
	"regexp"

	"github.com/gomponents/gontainer/pkg/syntax"
)

var (
	serviceNameRegex = regexp.MustCompile(`^` + syntax.ServiceNamePattern + `$`)
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
