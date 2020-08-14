package arguments

import (
	"fmt"
	"strings"

	"github.com/gomponents/gontainer-helpers/exporters"
	"github.com/gomponents/gontainer/pkg/dto"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/tokens"
)

type PatternResolver struct {
	tokenizer tokens.Tokenizer
	exporter  exporters.Exporter
	imports   imports.Imports
}

func NewPatternResolver(tokenizer tokens.Tokenizer, exporter exporters.Exporter, imports imports.Imports) *PatternResolver {
	return &PatternResolver{tokenizer: tokenizer, exporter: exporter, imports: imports}
}

func (p PatternResolver) Resolve(expr string) (dto.CompiledArg, error) {
	tkns, err := p.tokenizer.Tokenize(expr)

	if err != nil {
		return dto.CompiledArg{}, fmt.Errorf("cannot tokenize expression %s", expr)
	}

	if len(tkns) == 0 {
		return dto.CompiledArg{}, fmt.Errorf("unexpected error, tokenizer returned 0 tokens for expression `%s`", expr)
	}

	var dependsOn []string

	solveTokenCode := func(t tokens.Token) (string, error) {
		switch t.Kind {
		case tokens.TokenKindString:
			return p.exporter.Export(t.Raw)
		case tokens.TokenKindReference:
			runes := []rune(t.Raw)
			depID := string(runes[1 : len(runes)-1])
			dependsOn = append(dependsOn, depID)
			return fmt.Sprintf("result.MustGetParam(%+q)", depID), nil
		case tokens.TokenKindCode:
			return t.Code, nil
		default:
			return "", fmt.Errorf("unexpected TokenKind %v", t.Kind)
		}
	}

	if len(tkns) == 1 {
		t := tkns[0]
		code, err := solveTokenCode(t)
		if err != nil {
			return dto.CompiledArg{}, err
		}

		return dto.CompiledArg{
			Code:      code,
			DependsOn: dependsOn,
		}, nil
	}

	codeParts := make([]string, 0)
	for _, t := range tkns {
		tmp, solveErr := solveTokenCode(t)
		if solveErr != nil {
			return dto.CompiledArg{}, solveErr
		}

		tmp = p.imports.GetAlias("github.com/gomponents/gontainer-helpers/exporters") + `.MustToString(` + tmp + `)`
		codeParts = append(codeParts, tmp)
	}

	return dto.CompiledArg{
		Code:      strings.Join(codeParts, " + "),
		DependsOn: dependsOn,
	}, nil
}

func (p PatternResolver) Supports(string) bool {
	return true
}
