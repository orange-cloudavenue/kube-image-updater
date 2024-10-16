package rules

import (
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"

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

// ! semver-major rule

func (s *semverMajor) Evaluate() (matchWithRule bool, newTag string, err error) {
	x, err := semver.NewVersion(s.actualTag)
	if err != nil {
		return false, "", err
	}

	// sort tags in descending order
	sort.Sort(sort.Reverse(sort.StringSlice(s.tags)))

	for _, t := range s.tags {
		ac, err := semver.NewVersion(t)
		if err != nil {
			log.WithError(err).WithField("tag", t).Error("Error parsing actual tag")
			continue
		}

		v, err := semver.NewConstraint(fmt.Sprintf("^%s", x.IncMajor()))
		if err != nil {
			log.WithError(err).WithField("constraint", x.IncMajor()).Error("Error parsing constraint")
			continue
		}

		if v.Check(ac) {
			s.SetNewTag(t)
			return true, t, nil
		}
	}

	return false, "", nil
}

// ! semver-minor rule

func (s *semverMinor) Evaluate() (matchWithRule bool, newTag string, err error) {
	x, err := semver.NewVersion(s.actualTag)
	if err != nil {
		return false, "", err
	}

	// sort tags in descending order
	sort.Sort(sort.Reverse(sort.StringSlice(s.tags)))

	for _, t := range s.tags {
		ac, err := semver.NewVersion(t)
		if err != nil {
			log.WithError(err).WithField("tag", t).Error("Error parsing actual tag")
			continue
		}

		v, err := semver.NewConstraint(fmt.Sprintf("^%s", x.IncMinor()))
		if err != nil {
			log.WithError(err).WithField("constraint", x.IncMinor()).Error("Error parsing constraint")
			continue
		}

		if v.Check(ac) {
			s.SetNewTag(t)
			return true, t, nil
		}
	}

	return false, "", nil
}

// ! semver-patch rule

func (s *semverPatch) Evaluate() (matchWithRule bool, newTag string, err error) {
	x, err := semver.NewVersion(s.actualTag)
	if err != nil {
		return false, "", err
	}

	// sort tags in descending order
	sort.Sort(sort.Reverse(sort.StringSlice(s.tags)))

	for _, t := range s.tags {
		ac, err := semver.NewVersion(t)
		if err != nil {
			log.WithError(err).WithField("tag", t).Error("Error parsing actual tag")
			continue
		}

		v, err := semver.NewConstraint(fmt.Sprintf(">=%s <%s", x.IncPatch(), x.IncMinor()))
		if err != nil {
			log.WithError(err).WithField("constraint", x.IncPatch()).Error("Error parsing constraint")
			continue
		}

		if v.Check(ac) {
			s.SetNewTag(t)
			return true, t, nil
		}
	}

	return false, "", nil
}
