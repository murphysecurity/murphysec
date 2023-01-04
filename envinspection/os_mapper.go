package envinspection

import "regexp"

type rule struct {
	Pattern     *regexp.Regexp `json:"pattern"`
	Replacement string         `json:"replacement"`
}

func (r *rule) Replace(input string) (string, bool) {
	if r.Pattern.MatchString(input) {
		return r.Pattern.ReplaceAllString(input, r.Replacement), true
	}
	return "", false
}

var rules = []rule{
	{regexp.MustCompile(`.*ctyunos.*`), "centos:8"},
}

func processByRule(input string) (string, bool) {
	for _, r := range rules {
		s, ok := r.Replace(input)
		if ok {
			return s, true
		}
	}
	return input, false
}
