package arguments

import (
	"fmt"
	"regexp"

	"github.com/gomponents/gontainer/pkg/dto/compiled"
	"github.com/gomponents/gontainer/pkg/regex"
)

var (
	taggedRegex = regexp.MustCompile("^" + regex.TaggedArg + "$")
)

type TaggedResolver struct {
}

func NewTaggedResolver() *TaggedResolver {
	return &TaggedResolver{}
}

func (t TaggedResolver) Resolve(p interface{}) (compiled.Arg, error) {
	s, _ := p.(string)
	_, m := regex.Match(taggedRegex, s)

	return compiled.Arg{
		Code:          fmt.Sprintf("result.MustGetByTag(%+q)", m["tag"]),
		Raw:           s,
		DependsOnTags: []string{m["tag"]},
	}, nil
}

func (t TaggedResolver) Supports(p interface{}) bool {
	s, ok := p.(string)
	if !ok {
		return false
	}
	return taggedRegex.MatchString(s)
}
