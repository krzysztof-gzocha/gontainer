package tokens

import (
	"fmt"
	"github.com/gomponents/gontainer/pkg/regex"
	"regexp"

	"github.com/gomponents/gontainer/pkg/imports"
)

type TokenKind uint

const (
	TokenKindString TokenKind = iota
	TokenKindReference
	TokenKindCode
)

const (
	TokenDelimiter      = "%"
	RegexTokenReference = `^[a-zA-Z][a-zA-Z0-9_]*((\.)[a-zA-Z0-9_]+)*$`
	RegexSimpleFn       = `^(?P<fn>[a-zA-Z][a-zA-Z0-9_]*((\.)[a-zA-Z0-9_]+)*)\((?P<params>.*)\)$`
)

var (
	regexTokenRef = regexp.MustCompile(RegexTokenReference)
	regexSimpleFn = regexp.MustCompile(RegexSimpleFn)
)

type Token struct {
	Kind      TokenKind
	Raw       string
	DependsOn []string
	Code      string
}

type TokenFactoryStrategy interface {
	Supports(expr string) bool
	Create(expr string) (Token, error)
}

// toExpr removes surrounding delimiters
func toExpr(expr string) (string, bool) {
	runes := []rune(expr)
	if len(runes) < 2 {
		return "", false
	}

	if string(runes[0]) != TokenDelimiter || string(runes[len(runes)-1]) != TokenDelimiter {
		return "", false
	}

	return string(runes[1 : len(runes)-1]), true
}

// %%
type TokenPercentSign struct{}

func (t TokenPercentSign) Supports(expr string) bool {
	return expr == "%%"
}

func (t TokenPercentSign) Create(expr string) (Token, error) {
	return Token{
		Kind:      TokenKindString,
		Raw:       "%%",
		DependsOn: nil,
	}, nil
}

// %my.param%
type TokenReference struct{}

func (t TokenReference) Supports(s string) bool {
	expr, ok := toExpr(s)

	return ok && regexTokenRef.MatchString(expr)
}

func (t TokenReference) Create(s string) (Token, error) {
	ref, _ := toExpr(s)

	return Token{
		Kind:      TokenKindReference,
		Raw:       s,
		DependsOn: []string{ref},
	}, nil
}

type TokenString struct{}

func (t TokenString) Supports(expr string) bool {
	return true
}

func (t TokenString) Create(expr string) (Token, error) {
	return Token{
		Kind:      TokenKindString,
		Raw:       expr,
		DependsOn: nil,
	}, nil
}

// %env(ENV_VAR)%
type TokenSimpleFunction struct {
	imports  imports.Imports
	fn       string
	goImport string
	goFn     string
}

func NewTokenSimpleFunction(imports imports.Imports, fn string, goImport string, goFn string) *TokenSimpleFunction {
	return &TokenSimpleFunction{imports: imports, fn: fn, goImport: goImport, goFn: goFn}
}

func (t TokenSimpleFunction) Supports(expr string) bool {
	e, ok := toExpr(expr)
	if !ok {
		return false
	}

	ok, m := regex.Match(regexSimpleFn, e)
	return ok && m["fn"] == t.fn
}

// todo it won't work without alias
func (t TokenSimpleFunction) Create(expr string) (Token, error) {
	e, _ := toExpr(expr)
	_, m := regex.Match(regexSimpleFn, e)
	return Token{
		Kind: TokenKindCode,
		Raw:  expr,
		Code: fmt.Sprintf("%s.%s(%s)", t.imports.GetAlias(t.goImport), t.goFn, m["params"]),
	}, nil
}
