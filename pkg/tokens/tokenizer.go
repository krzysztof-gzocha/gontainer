package tokens

import (
	"fmt"
)

type Tokenizer interface {
	Tokenize(pattern string) ([]Token, error)
}

type PatternTokenizer struct {
	strategies []TokenFactoryStrategy
}

func NewPatternTokenizer(strategies []TokenFactoryStrategy) *PatternTokenizer {
	return &PatternTokenizer{strategies: strategies}
}

func (s PatternTokenizer) Tokenize(pattern string) ([]Token, error) {
	result := make([]Token, 0)

	opened := false
	buff := ""

	// todo use runes
	for _, c := range pattern {
		ch := string(c)
		if ch == TokenDelimiter {
			if opened {
				if t, err := s.createToken(buff + TokenDelimiter); err == nil {
					result = append(result, t)
					buff = ""
					opened = false
					continue
				} else {
					return nil, err
				}
			}

			if buff != "" {
				if t, err := s.createToken(buff); err == nil {
					result = append(result, t)
				} else {
					return nil, err
				}
			}
			buff = TokenDelimiter
			opened = true
			continue
		}

		buff += ch
	}

	if buff != "" {
		if t, err := s.createToken(buff); err == nil {
			result = append(result, t)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func (s PatternTokenizer) createToken(expr string) (Token, error) {
	for _, strategy := range s.strategies {
		if strategy.Supports(expr) {
			t, err := strategy.Create(expr)
			if err != nil {
				err = fmt.Errorf("cannot create token `%s`: %s", expr, err.Error())
			}
			return t, err
		}
	}

	return Token{}, fmt.Errorf("invalid token `%s`", expr)
}
