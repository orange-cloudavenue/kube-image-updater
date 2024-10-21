package rules

import (
	"go/version"
	"sort"
)

var _ RuleInterface = &always{}

type (
	always struct {
		rule
	}
)

func init() {
	register(Always, &always{})
}

// ! always rule

func (a *always) Evaluate() (matchWithRule bool, newTag string, err error) {
	if len(a.tags) == 0 {
		return false, "", nil
	}

	// Sort the tags
	sort.Slice(a.tags, func(i, j int) bool {
		return version.Compare(a.tags[i], a.tags[j]) == -1
	})

	// Return the last tag if it is not the actual tag
	if a.tags[len(a.tags)-1] != a.actualTag {
		return true, a.tags[len(a.tags)-1], nil
	}

	return false, "", nil
}
