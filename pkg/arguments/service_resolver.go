package arguments

import (
	"strings"
)

type ServiceResolver struct{}

// todo syntax validation
// check whether service exists
// @serviceName.(foo/bar/foobar.MyType)
// @service
func (s ServiceResolver) Resolve(expr string) (Argument, error) {
	parts := strings.Split(expr, ".(")

	runeName := []rune(parts[0])
	name := string(runeName[1:])

	if len(parts) == 1 {
		return Argument{
			Kind: ArgumentKindService,
			ServiceMetadata: ServiceMetadata{
				ID:          name,
				Import:      "",
				Type:        "",
				PointerType: false,
			},
		}, nil
	}

	runeType := []rune(parts[1])
	fullType := string(runeType[:len(runeType)-1])
	typeParts := strings.Split(fullType, ".")

	import_ := strings.Join(typeParts[:len(typeParts)-1], ".")
	type_ := typeParts[len(typeParts)-1]
	pointer := false

	if []rune(import_)[0] == '*' {
		import_ = import_[1:]
		pointer = true
	}

	return Argument{
		Kind: ArgumentKindService,
		ServiceMetadata: ServiceMetadata{
			ID:          name,
			Import:      import_,
			Type:        type_,
			PointerType: pointer,
		},
	}, nil
}

func (s ServiceResolver) Supports(expr string) bool {
	return []rune(expr)[0] == '@'
}
