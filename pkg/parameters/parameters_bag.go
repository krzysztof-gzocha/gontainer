package parameters

import (
	"fmt"
	"sort"
	"strings"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/tokens"
)

type Imports interface {
	GetAlias(string) string
}

type Exporter interface {
	Export(interface{}) (string, error)
}

type bagParam struct {
	resolved bool
	param    compiled.Param
}

type bagParams map[string]bagParam

func (b bagParams) SortedKeys() []string {
	var r []string
	for k, _ := range b {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}

type BagFactory struct {
	tokenizer tokens.Tokenizer
	exporter  Exporter
	imports   Imports
}

func NewBagFactory(tokenizer tokens.Tokenizer, exporter Exporter, imports Imports) *BagFactory {
	return &BagFactory{tokenizer: tokenizer, exporter: exporter, imports: imports}
}

func (s BagFactory) Create(params map[string]interface{}) (map[string]compiled.Param, error) {
	bag := s.createWrappedBag(params)

	if err := s.solveNonStrings(bag); err != nil {
		return nil, err
	}

	if err := s.solveStrings(bag); err != nil {
		return nil, err
	}

	return s.unwrapBag(bag), nil
}

func (s BagFactory) solveNonStrings(bag bagParams) error {
	for _, n := range bag.SortedKeys() {
		p := bag[n]
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

func (s BagFactory) solveStrings(bag bagParams) error {
	for _, id := range bag.SortedKeys() {
		p := bag[id]
		if _, ok := p.param.Raw.(string); !ok {
			continue
		}

		if err := s.solveString(id, bag, newDependenciesBag()); err != nil {
			return err
		}
	}

	return nil
}

func (s BagFactory) solveString(id string, bag bagParams, deps dependenciesBag) error {
	if _, ok := bag[id]; !ok {
		return fmt.Errorf("parameter `%s` does not exist", id)
	}

	if bag[id].resolved {
		return nil
	}

	if deps.Has(id) {
		return fmt.Errorf("cannot solve param `%s`, circular dependencies: %s", id, deps.String())
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
			// todo make "result" injectable
			return fmt.Sprintf("result.MustGetParam(%+q)", depID), nil
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

func (BagFactory) createWrappedBag(params map[string]interface{}) bagParams {
	result := make(bagParams)
	for n, p := range params {
		result[n] = bagParam{
			resolved: false,
			param:    compiled.Param{Raw: p},
		}
	}

	return result
}

func (BagFactory) unwrapBag(bag bagParams) map[string]compiled.Param {
	result := make(map[string]compiled.Param)
	for n, p := range bag {
		result[n] = p.param
	}
	return result
}
