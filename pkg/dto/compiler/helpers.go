package compiler

import (
	"strings"
)

func sanitizeImport(i string) string {
	return strings.Trim(i, `"`)
}
