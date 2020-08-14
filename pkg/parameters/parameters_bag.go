package parameters

import (
	"fmt"
	"strings"

	"github.com/gomponents/gontainer-helpers/exporters"
	"github.com/gomponents/gontainer/pkg/imports"
	"github.com/gomponents/gontainer/pkg/tokens"
)

type bagParam struct {
	resolved bool
	param    ResolvedParam
}

type bagParams map[string]bagParam

type BagFactory interface {
	Create(RawParameters) (ResolvedParams, error)
}

type SimpleBagFactory struct {
	tokenizer tokens.Tokenizer
	exporter  exporters.Exporter
	imports   imports.Imports
}

func NewSimpleBagFactory(tokenizer tokens.Tokenizer, exporter exporters.Exporter, imports imports.Imports) *SimpleBagFactory {
	return &SimpleBagFactory{tokenizer: tokenizer, exporter: exporter, imports: imports}
}

func (s SimpleBagFactory) Create(params RawParameters) (ResolvedParams, error) {
	bag := s.createWrappedBag(params)

	if err := s.solveNonStrings(bag); err != nil {
		return nil, err
	}

	if err := s.solveStrings(bag); err != nil {
		return nil, err
	}

	return s.unwrapBag(bag), nil
}

func (s SimpleBagFactory) solveNonStrings(bag bagParams) error {
	for n, p := range bag {
		if _, ok := p.param.Raw.(string); ok {
			continue
		}

		code, err := s.exporter.Export(p.param.Raw)
		if err != nil {
			return fmt.Errorf("cannot export param: `%s`: %s", n, err.Error())
		}

		p.param.Code = code
		p.resolved = true

		bag[n] = p
	}

	return nil
}

func (s SimpleBagFactory) solveStrings(bag bagParams) error {
	for id, p := range bag {
		if _, ok := p.param.Raw.(string); !ok {
			continue
		}

		if err := s.solveString(id, bag, newDependenciesBag()); err != nil {
			return err
		}
	}

	return nil
}

func (s SimpleBagFactory) solveString(id string, bag bagParams, deps dependenciesBag) error {
	if _, ok := bag[id]; !ok {
		return fmt.Errorf("parameter `%s` does not exist", id)
	}

	if bag[id].resolved {
		return nil
	}

	if deps.Has(id) {
		return fmt.Errorf("cannot solve param `%s`, circular dependencies: %s", id, deps.ToString())
	}

	pattern := bag[id].param.Raw.(string)
	tkns, tokenizerErr := s.tokenizer.Tokenize(pattern)
	if tokenizerErr != nil {
		return fmt.Errorf("cannot solve param `%s`: %s", id, tokenizerErr.Error())
	}

	solveDep := func(depID string) error {
		depsClone := deps.Clone()
		depsClone.Append(id)

		return s.solveString(depID, bag, depsClone)
	}

	// todo refactor
	solveTokenCode := func(t tokens.Token) (string, error) {
		switch t.Kind {
		case tokens.TokenKindString:
			return s.exporter.Export(t.Raw)
		case tokens.TokenKindReference:
			runes := []rune(t.Raw)
			depID := string(runes[1 : len(runes)-1])
			if err := solveDep(depID); err != nil {
				return "", err
			}
			return fmt.Sprintf("result.MustGetParam(%+q)", depID), nil
			// todo
			//return bag[depID].param.Code, nil
		case tokens.TokenKindCode:
			return t.Code, nil
		default:
			return "", fmt.Errorf("unexpected TokenKind %v", t.Kind)
		}
	}

	if len(tkns) == 0 {
		return fmt.Errorf("unexpected error, tokenizer returned 0 tokens for parameter `%s`", id)
	}

	if len(tkns) == 1 {
		t := tkns[0]
		code, solveErr := solveTokenCode(t)
		if solveErr != nil {
			return solveErr
		}

		cp := bag[id]
		cp.param.Code = code
		cp.resolved = true
		bag[id] = cp

		return nil
	}

	codeParts := make([]string, 0)
	for _, t := range tkns {
		tmp, solveErr := solveTokenCode(t)
		if solveErr != nil {
			return solveErr
		}

		tmp = s.imports.GetAlias("github.com/gomponents/gontainer-helpers/exporters") + `.MustToString(` + tmp + `)`
		codeParts = append(codeParts, tmp)
	}

	cp := bag[id]
	cp.resolved = true
	cp.param.Code = strings.Join(codeParts, " + ")
	bag[id] = cp

	return nil
}

func (SimpleBagFactory) createWrappedBag(params RawParameters) bagParams {
	result := make(bagParams)
	for n, p := range params {
		result[n] = bagParam{
			resolved: false,
			param:    ResolvedParam{Raw: p},
		}
	}

	return result
}

func (SimpleBagFactory) unwrapBag(bag bagParams) ResolvedParams {
	result := make(ResolvedParams)
	for n, p := range bag {
		result[n] = p.param
	}
	return result
}
