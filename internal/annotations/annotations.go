package annotations

import (
	"context"
	"crypto/md5" //nolint:gosec
	"encoding/json"
	"fmt"
	"strings"
)

type (
	annotations map[string]string

	aChan chan struct {
		key   AnnotationKey
		value string
	}

	Annotation struct {
		ctx context.Context
		annotations
	}

	Action struct {
		aChan aChan
		value string
	}

	Tag struct {
		aChan aChan
		value string
	}

	CheckSum struct {
		aChan aChan
		value string
	}

	AnnotationKey string
)

// AnnotationKey is the key used to store the image in the annotation
var (
	KeyAction   AnnotationKey = "kimup.cloudavenue.io" + "/action"
	KeyTag      AnnotationKey = "kimup.cloudavenue.io" + "/tag"
	KeyCheckSum AnnotationKey = "kimup.cloudavenue.io" + "/checksum"
)

type (
	KubeAnnotationInterface interface {
		GetAnnotations() map[string]string
	}
)

func New(ctx context.Context, object KubeAnnotationInterface) Annotation {
	return Annotation{
		ctx: ctx,
		// src:         object,
		annotations: object.GetAnnotations(),
	}
}

// * Global

func (a *Annotation) Remove(key AnnotationKey) {
	delete(a.annotations, string(key))
}

// * Action

type AActionKey string

const (

	// Add Actions here

	// Action Refresh
	ActionRefresh AActionKey = "refresh"
	ActionDelete  AActionKey = "delete"
)

func (a *Annotation) Action() (ac *Action) {
	ac = &Action{
		aChan: make(aChan),
	}

	if v, ok := a.annotations[string(KeyAction)]; ok {
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

func (a *Action) Is(action AActionKey) bool {
	return strings.EqualFold(a.value, string(action))
}

func (a *Action) IsNull() bool {
	return a.value == ""
}

func (a *Action) Get() AActionKey {
	return AActionKey(a.value)
}

func (a *Action) Set(action AActionKey) {
	a.aChan.Send(KeyAction, string(action))
}

// * Tag

func (a *Annotation) Tag() (at *Tag) {
	at = &Tag{
		aChan: make(aChan),
	}

	if v, ok := a.annotations[string(KeyTag)]; ok {
		at.value = v
	}

	go func() {
		for {
			select {
			case x := <-at.aChan:
				a.annotations[string(x.key)] = x.value
			case <-a.ctx.Done():
				return
			}
		}
	}()

	return at
}

func (a *Tag) Get() string {
	return a.value
}

func (a *Tag) Set(tag string) {
	a.aChan.Send(KeyTag, tag)
}

func (a *Tag) IsNull() bool {
	return a.value == ""
}

// * CheckSum

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

// * Generic funcs

func (aC aChan) Send(key AnnotationKey, value string) {
	aC <- struct {
		key   AnnotationKey
		value string
	}{key, value}
}
