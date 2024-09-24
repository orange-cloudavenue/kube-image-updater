package rules

import (
	log "github.com/sirupsen/logrus"
)

type (
	RuleInterface interface {
		Init(actualTag string, tagsAvailable []string)
		Evaluate() (matchWithRule bool, newTag string, err error)
		GetNewTag() string
	}

	Rules map[Name]RuleInterface
	Name  string
)

var rules = make(Rules)

const (
	SemverMajor Name = "semver-major"
	SemverMinor Name = "semver-minor"
	SemverPatch Name = "semver-patch"
	// SemverPreRelease Name = "semver-prerelease"
	Regex Name = "regex"
)

func RegisterRule(name Name, rule RuleInterface) {
	log.Infof("Registering rule %s", name)
	rules[name] = rule
}

func GetRule(name Name) RuleInterface {
	return rules[name]
}
