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

func RegisterRule(name Name, rule RuleInterface) {
	log.Infof("Registering rule %s", name)
	rules[name] = rule
}

func GetRule(name Name) RuleInterface {
	return rules[name]
}

// * Generic func set tags for all semver rules
func (r *rule) Init(actualTag string, tagsAvailable []string, value string) {
	r.actualTag = actualTag
	r.tags = tagsAvailable
	r.value = value
}

func (r *rule) GetNewTag() string {
	return r.newTag
}
