package regex

import (
	"regexp"
)

func Match(r *regexp.Regexp, s string) (bool, map[string]string) {
	if !r.MatchString(s) {
		return false, nil
	}

	match := r.FindStringSubmatch(s)
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return true, result
}
