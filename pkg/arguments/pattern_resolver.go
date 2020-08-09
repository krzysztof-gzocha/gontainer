package arguments

import (
	"fmt"
	"strings"

	"github.com/gomponents/gontainer/pkg/dto"
	"github.com/gomponents/gontainer/pkg/exporters"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/parameters"
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

func (p PatternResolver) Resolve(expr string, params parameters.ResolvedParams) (dto.CompiledArg, error) {
	tkns, err := p.tokenizer.Tokenize(expr)

	if err != nil {
		return dto.CompiledArg{}, fmt.Errorf("cannot tokenize expression %s", expr)
	}

	if len(tkns) == 0 {
		return dto.CompiledArg{}, fmt.Errorf("unexpected error, tokenizer returned 0 tokens for expression `%s`", expr)
	}

	solveTokenCode := func(t tokens.Token) (string, error) {
		switch t.Kind {
		case tokens.TokenKindString:
			return p.exporter.Export(t.Raw)
		case tokens.TokenKindReference:
			runes := []rune(t.Raw)
			depID := string(runes[1 : len(runes)-1])
			param, ok := params[depID]
			if !ok {
				return "", fmt.Errorf("parameter `%s` does not exist", depID)
			}
			return param.Code, nil
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

		return dto.CompiledArg{Code: code}, nil
	}

	codeParts := make([]string, 0)
	for _, t := range tkns {
		tmp, solveErr := solveTokenCode(t)
		if solveErr != nil {
			return dto.CompiledArg{}, solveErr
		}

		tmp = p.imports.GetAlias("github.com/gomponents/gontainer/pkg/std") + `.MustConvertToString(` + tmp + `)`
		codeParts = append(codeParts, tmp)
	}

	return dto.CompiledArg{Code: strings.Join(codeParts, " + ")}, nil
}

func (p PatternResolver) Supports(string) bool {
	return true
}
