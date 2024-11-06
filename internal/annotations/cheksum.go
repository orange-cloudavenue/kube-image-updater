package annotations

import (
	"crypto/md5" //nolint:gosec
	"encoding/json"
	"fmt"
)

// * CheckSum

type (
	CheckSum struct {
		aChan aChan
		value string
	}
)

func (a *Annotation) CheckSum() (ac *CheckSum) {
	ac = &CheckSum{
		aChan: make(aChan),
	}

	if v, ok := a.annotations[string(KeyCheckSum)]; ok {
		ac.value = v
	}

	go func() {
		for {
			select {
			case x := <-ac.aChan:
				a.annotations[string(x.key)] = x.value
			case <-a.ctx.Done():
				return
			}
		}
	}()

	return ac
}

func (a *CheckSum) Get() string {
	return a.value
}

func (a *CheckSum) Set(object interface{}) error {
	x, err := a.computeChecksum(object)
	if err != nil {
		return err
	}

	a.aChan.Send(KeyCheckSum, x)
	return nil
}

func (a *CheckSum) IsNull() bool {
	return a.value == ""
}

func (a *CheckSum) computeChecksum(object interface{}) (string, error) {
	x, err := json.Marshal(object)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", (md5.Sum(x))), nil //nolint:gosec
}

// IsEqual compares the checksum of the object with the one stored in the annotation
func (a *CheckSum) IsEqual(object interface{}) (bool, error) {
	if a.IsNull() {
		return false, nil
	}

	x, err := a.computeChecksum(object)
	if err != nil {
		return false, err
	}

	return a.value == x, nil
}
