package rules

import (
	"fmt"
	"sort"

	"github.com/shipengqi/vc"

	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
)

var (
	_ RuleInterface = &semverMajor{}
	_ RuleInterface = &semverMinor{}
	_ RuleInterface = &semverPatch{}
)

type (
	semverMajor struct {
		rule
	}

	semverMinor struct {
		rule
	}

	semverPatch struct {
		rule
	}
)

func init() {
	register(SemverMajor, &semverMajor{})
	register(SemverMinor, &semverMinor{})
	register(SemverPatch, &semverPatch{})
}

var funcParseSemVer = func(s string) (vc.Comparable, error) {
	return vc.NewSemverStr(s)
}

// ! semver-major rule

func (s *semverMajor) Evaluate() (matchWithRule bool, newTag string, err error) {
	x, err := vc.NewSemverStr(s.actualTag)
	if err != nil {
		return false, "", err
	}

	// sort tags in descending order
	sort.Sort(sort.Reverse(sort.StringSlice(s.tags)))

	for _, t := range s.tags {
		// Original x = 1.0.0
		// >=2.0.0
		v, err := vc.NewConstraint(fmt.Sprintf(">=%s", x.IncMajor()), funcParseSemVer)
		if err != nil {
			log.WithError(err).WithField("constraint", x.IncMinor()).Error("Error parsing constraint")
			continue
		}

		if ok, _ := v.CheckString(t); ok {
			s.SetNewTag(t)
			return true, t, nil
		}
	}

	return false, "", nil
}

// ! semver-minor rule

func (s *semverMinor) Evaluate() (matchWithRule bool, newTag string, err error) {
	x, err := vc.NewSemverStr(s.actualTag)
	if err != nil {
		return false, "", err
	}

	// sort tags in descending order
	sort.Sort(sort.Reverse(sort.StringSlice(s.tags)))

	for _, t := range s.tags {
		// Original x = 1.0.0
		// >=1.1.0 <2
		v, err := vc.NewConstraint(fmt.Sprintf(">=%s <%s", x.IncMinor(), x.IncMajor()), funcParseSemVer)
		if err != nil {
			log.WithError(err).WithField("constraint", x.IncMinor()).Error("Error parsing constraint")
			continue
		}

		if ok, _ := v.CheckString(t); ok {
			s.SetNewTag(t)
			return true, t, nil
		}
	}

	return false, "", nil
}

// ! semver-patch rule

func (s *semverPatch) Evaluate() (matchWithRule bool, newTag string, err error) {
	x, err := vc.NewSemverStr(s.actualTag)
	if err != nil {
		return false, "", err
	}

	// sort tags in descending order
	sort.Sort(sort.Reverse(sort.StringSlice(s.tags)))

	for _, t := range s.tags {
		// Original x = 1.0.0
		// >=1.0.1 <1.1.0
		v, err := vc.NewConstraint(fmt.Sprintf(">=%s <%s", x.IncPatch(), x.IncMinor()), funcParseSemVer)
		if err != nil {
			log.WithError(err).WithField("constraint", x.IncMinor()).Error("Error parsing constraint")
			continue
		}

		if ok, _ := v.CheckString(t); ok {
			s.SetNewTag(t)
			return true, t, nil
		}
	}

	return false, "", nil
}
