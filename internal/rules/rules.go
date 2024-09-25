package rules

import (
	log "github.com/sirupsen/logrus"
)

type (
	RuleInterface interface {
		Init(actualTag string, tagsAvailable []string, value string)
		Evaluate() (matchWithRule bool, newTag string, err error)
		GetNewTag() string
	}

	Rules map[Name]RuleInterface
	Name  string

	rule struct {
		actualTag string
		newTag    string
		tags      []string
		value     string
	}
)

var rules = make(Rules)

const (
	SemverMajor Name = "semver-major"
	SemverMinor Name = "semver-minor"
	SemverPatch Name = "semver-patch"
	Regex       Name = "regex"
)

func register(name Name, rule RuleInterface) {
	log.Infof("Registering rule %s", name)
	rules[name] = rule
}

// GetRule retrieves a RuleInterface based on the provided name.
// It takes a Name type as an argument and returns the corresponding RuleInterface
// from the rules map. If the name does not exist in the map, the behavior is
// dependent on the implementation of the rules map.
//
// Parameters:
//   - name: The name of the rule to retrieve.
//
// Returns:
//   - RuleInterface: The rule associated with the given name.
func GetRule(name Name) (RuleInterface, error) {
	if _, ok := rules[name]; !ok {
		return nil, ErrRuleNotFound
	}
	return rules[name], nil
}

// GetRuleWithUntypedName retrieves a RuleInterface based on the provided name.
// It takes a string as an argument and returns the corresponding RuleInterface
// from the rules map. If the name does not exist in the map, the behavior is
// dependent on the implementation of the rules map.
func GetRuleWithUntypedName(name string) (RuleInterface, error) {
	n, err := ParseRuleName(name)
	if err != nil {
		return nil, err
	}
	return GetRule(n)
}

// ParseRuleName takes a rule name as a string and checks if it exists in the predefined rules.
// If the rule is found, it returns the corresponding Name and a nil error.
// If the rule is not found, it returns an empty Name and an error indicating that the rule was not found.
//
// Parameters:
//   - name: A string representing the name of the rule to be parsed.
//
// Returns:
//   - Name: The corresponding Name if found.
//   - error: An error if the rule is not found.
func ParseRuleName(name string) (Name, error) {
	for k := range rules {
		if string(k) == name {
			return k, nil
		}
	}
	return "", ErrRuleNotFound
}

// Init initializes the rule with the provided actual tag, available tags, and value.
//
// Parameters:
//   - actualTag: The current tag of the image.
//   - tagsAvailable: A slice of tags that are available for the image.
//   - value: A string value associated with the rule.
func (r *rule) Init(actualTag string, tagsAvailable []string, value string) {
	r.actualTag = actualTag
	r.tags = tagsAvailable
	r.value = value
}

// GetNewTag returns the new tag associated with the rule.
// It retrieves the value of the newTag field from the rule instance.
func (r *rule) GetNewTag() string {
	return r.newTag
}

// String returns the string representation of the rule name.
func (r Name) String() string {
	return string(r)
}
