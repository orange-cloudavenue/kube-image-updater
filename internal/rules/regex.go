package rules

import (
	"fmt"
	"regexp"
	"sort"
)

var _ RuleInterface = &regex{}

type (
	regex struct {
		rule
	}
)

func init() {
	register(Regex, &regex{})
}

// ! regex rule

func (r *regex) Evaluate() (matchWithRule bool, newTag string, err error) {
	// sort tags in descending order
	sort.Sort(sort.Reverse(sort.StringSlice(r.tags)))

	if r.value == "" {
		return false, "", fmt.Errorf("regex value is empty")
	}

	re, err := regexp.Compile(r.value)
	if err != nil {
		return false, "", err
	}

	for _, t := range r.tags {
		if t == r.actualTag {
			continue
		}
		if re.MatchString(t) {
			return true, t, nil
		}
	}

	return false, "", nil
}
