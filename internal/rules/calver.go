package rules

import (
	"fmt"

	"github.com/shipengqi/vc"

	"github.com/orange-cloudavenue/kube-image-updater/internal/log"
)

type (
	// calverMajor - The first number in the version.
	calverMajor struct {
		rule
	}

	// calverMinor - The second number in the version.
	calverMinor struct {
		rule
	}

	// calverPatch - The third and usually final number in the version. Sometimes referred to as the "micro" segment.
	calverPatch struct {
		rule
	}

	// calverPrerelease - The prerelease is an optional part of the version.
	calverPrerelease struct {
		rule
	}
)

// New a function to generate a Comparable instance.
var funcParseCalver = func(s string) (vc.Comparable, error) {
	return vc.NewCalVerStr(s)
}

func init() {
	register(CalverMajor, &calverMajor{})
	register(CalverMinor, &calverMinor{})
	register(CalverPatch, &calverPatch{})
	register(CalverPrerelease, &calverPrerelease{})
}

func (c *calverMajor) Evaluate() (matchWithRule bool, newTag string, err error) {
	actualCV, err := vc.NewCalVerStr(c.actualTag)
	if err != nil {
		log.WithError(err).WithField("tag", c.actualTag).Error("Error parsing actual tag")
		return false, "", err
	}

	for _, t := range c.tags {
		cv, err := vc.NewCalVerStr(t)
		if err != nil {
			log.WithError(err).WithField("tag", t).Error("Error parsing tag")
			continue
		}

		// Contains no prerelease
		if cv.Prerelease() == "" {
			// Create a constraint (e.g. 2024.0.0: >=2025.0.0)
			constraint, err := vc.NewConstraint(fmt.Sprintf(">=%s", actualCV.IncMajor()), funcParseCalver)
			if err != nil {
				log.WithError(err).WithField("constraint", t).Error("Error parsing constraint")
				continue
			}

			// Check if the constraint is satisfied
			if constraint.Check(cv) {
				c.SetNewTag(t)
				return true, t, nil
			}
		}
	}

	return false, "", nil
}

func (c *calverMinor) Evaluate() (matchWithRule bool, newTag string, err error) {
	actualCV, err := vc.NewCalVerStr(c.actualTag)
	if err != nil {
		log.WithError(err).WithField("tag", c.actualTag).Error("Error parsing actual tag")
		return false, "", err
	}

	for _, t := range c.tags {
		cv, err := vc.NewCalVerStr(t)
		if err != nil {
			log.WithError(err).WithField("tag", t).Error("Error parsing tag")
			continue
		}

		// Contains no prerelease
		if cv.Prerelease() == "" {
			// Create a constraint (e.g. 2024.0.0: >=2024.1.0 <2025.0.0)
			constraint, err := vc.NewConstraint(fmt.Sprintf(">=%s <%s", actualCV.IncMinor(), actualCV.IncMajor()), funcParseCalver)
			if err != nil {
				log.WithError(err).WithField("constraint", t).Error("Error parsing constraint")
				continue
			}

			// Check if the constraint is satisfied
			if constraint.Check(cv) {
				c.SetNewTag(t)
				return true, t, nil
			}
		}
	}

	return false, "", nil
}

func (c *calverPatch) Evaluate() (matchWithRule bool, newTag string, err error) {
	actualCV, err := vc.NewCalVerStr(c.actualTag)
	if err != nil {
		log.WithError(err).WithField("tag", c.actualTag).Error("Error parsing actual tag")
		return false, "", err
	}

	for _, t := range c.tags {
		cv, err := vc.NewCalVerStr(t)
		if err != nil {
			log.WithError(err).WithField("tag", t).Error("Error parsing tag")
			continue
		}

		// Contains no prerelease
		if cv.Prerelease() == "" {
			// Create a constraint (e.g. 2024.0.0: >=2024.0.1 <2024.1.0)
			constraint, err := vc.NewConstraint(fmt.Sprintf(">=%s <%s", actualCV.IncPatch(), actualCV.IncMinor()), funcParseCalver)
			if err != nil {
				log.WithError(err).WithField("constraint", t).Error("Error parsing constraint")
				continue
			}

			// Check if the constraint is satisfied
			if constraint.Check(cv) {
				c.SetNewTag(t)
				return true, t, nil
			}
		}
	}

	return false, "", nil
}

// ***** calverModifier is not used in the current implementation *****
// func (c *calverModifier) Evaluate() (matchWithRule bool, newTag string, err error) {
// 	actualCV, err := vc.NewCalVerStr(c.actualTag)
// 	if err != nil {
// 		log.WithError(err).WithField("tag", c.actualTag).Error("Error parsing actual tag")
// 		return false, "", err
// 	}

// 	for _, t := range c.tags {
// 		cv, err := vc.NewCalVerStr(t)
// 		if err != nil {
// 			log.WithError(err).WithField("tag", t).Error("Error parsing tag")
// 			continue
// 		}

// 		// Contains a modifier
// 		if cv.Prerelease() != "" {
// 			// Create a constraint (e.g. 2024.0.0-dev: >=2024.0.1-dev)
// 			constraint, err := vc.NewConstraint(fmt.Sprintf(">=%s", actualCV.IncPatch()), funcParseCalver)
// 			if err != nil {
// 				log.WithError(err).WithField("constraint", t).Error("Error parsing constraint")
// 				continue
// 			}

// 			// Check if the constraint is satisfied
// 			if constraint.Check(cv) {
// 				c.SetNewTag(t)
// 				return true, t, nil
// 			}
// 		}
// 	}

// 	return false, "", nil
// }

func (c *calverPrerelease) Evaluate() (matchWithRule bool, newTag string, err error) {
	actualCV, err := vc.NewCalVerStr(c.actualTag)
	if err != nil {
		log.WithError(err).WithField("tag", c.actualTag).Error("Error parsing actual tag")
		return false, "", err
	}

	for _, t := range c.tags {
		cv, err := vc.NewCalVerStr(t)
		if err != nil {
			log.WithError(err).WithField("tag", t).Error("Error parsing tag")
			continue
		}

		// Contains a prerelease
		if cv.Prerelease() != "" {
			// Create a constraint (e.g. 2024.0.0-dev.0: >2024.0.0-dev.0)
			constraint, err := vc.NewConstraint(fmt.Sprintf(">%s-%s", actualCV.Version(), actualCV.Prerelease()), funcParseCalver)
			if err != nil {
				log.WithError(err).WithField("constraint", t).Error("Error parsing constraint")
				continue
			}

			// Check if the constraint is satisfied
			if constraint.Check(cv) {
				c.SetNewTag(t)
				return true, t, nil
			}
		}
	}

	return false, "", nil
}
