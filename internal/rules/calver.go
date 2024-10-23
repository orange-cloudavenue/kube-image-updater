package rules

import (
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
)

func init() {
	register(CalverMajor, &calverMajor{})
	register(CalverMinor, &calverMinor{})
	register(CalverPatch, &calverPatch{})
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

		if cv.Major() > actualCV.Major() {
			return true, t, nil
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

		if cv.Minor() > actualCV.Minor() && cv.Major() == actualCV.Major() {
			return true, t, nil
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

		if cv.Patch() > actualCV.Patch() && cv.Minor() == actualCV.Minor() && cv.Major() == actualCV.Major() {
			return true, t, nil
		}
	}

	return false, "", nil
}
