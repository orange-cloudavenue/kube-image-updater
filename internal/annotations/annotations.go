package annotations

import (
	"context"
	"crypto/md5" //nolint:gosec
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/orange-cloudavenue/kube-image-updater/internal/patch"
	"github.com/orange-cloudavenue/kube-image-updater/internal/utils"
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

	Enabled struct {
		aChan aChan
		value bool
	}

	MapImage struct {
		aChan aChan
		value map[string]string
	}

	MapContainer struct {
		aChan aChan
		value map[string]string
	}

	AnnotationKey string
)

// AnnotationKey is the key used to store the image in the annotation
var (
	KeyAction       AnnotationKey = "kimup.cloudavenue.io" + "/action"
	KeyTag          AnnotationKey = "kimup.cloudavenue.io" + "/tag"
	KeyCheckSum     AnnotationKey = "kimup.cloudavenue.io" + "/checksum"
	KeyEnabled      AnnotationKey = "kimup.cloudavenue.io" + "/enabled"
	KeyMapImage     AnnotationKey = "kimup.cloudavenue.io" + "/image"
	KeyMapContainer AnnotationKey = "kimup.cloudavenue.io" + "/container"
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
	ActionReload  AActionKey = "reload"
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

// * Enabled

func (a *Annotation) Enabled() Enabled {
	ae := Enabled{
		aChan: make(aChan),
	}

	if v, ok := a.annotations[string(KeyEnabled)]; ok {
		boolValue, _ := strconv.ParseBool(v)
		ae.value = boolValue
	}

	go func() {
		for {
			select {
			case x := <-ae.aChan:
				a.annotations[string(x.key)] = x.value
			case <-a.ctx.Done():
				return
			}
		}
	}()

	return ae
}

func (a Enabled) Get() bool {
	return a.value
}

func (a Enabled) Set(enabled bool) {
	a.aChan.Send(KeyEnabled, strconv.FormatBool(enabled))
}

// * Images

func (a *Annotation) Images() MapImage {
	ai := MapImage{
		aChan: make(aChan),
		value: make(map[string]string),
	}

	for k, v := range a.annotations {
		if strings.HasPrefix(k, string(KeyMapImage)) {
			ai.value[strings.TrimPrefix(k, string(KeyMapImage)+"/")] = v
		}
	}

	go func() {
		for {
			select {
			case x := <-ai.aChan:
				a.annotations[string(KeyMapImage)+"/"+string(x.key)] = x.value
			case <-a.ctx.Done():
				return
			}
		}
	}()

	return ai
}

func (a MapImage) Get(kubernetesImageName string) (crdImage string, err error) {
	if v, ok := a.value[kubernetesImageName]; ok {
		return v, nil
	}

	return "", fmt.Errorf("image %s not found", kubernetesImageName)
}

// Set associates a Kubernetes image name with a custom resource definition (CRD) image.
// It sends the mapping to the aChan channel using the specified AnnotationKey.
//
// Parameters:
//   - kubernetesImageName: The name of the Kubernetes image to be set.
//   - crdImage: The corresponding CRD image that is associated with the Kubernetes image.
func (a MapImage) Set(kubernetesImageName, crdImage string) {
	a.aChan.Send(AnnotationKey(kubernetesImageName), crdImage)
}

func (a MapImage) IsNull() bool {
	return len(a.value) == 0
}

// * Containers

func (a *Annotation) Containers() MapContainer {
	ai := MapContainer{
		aChan: make(aChan),
		value: make(map[string]string),
	}

	for k, v := range a.annotations {
		if strings.HasPrefix(k, string(KeyMapContainer)) {
			ai.value[strings.TrimPrefix(k, string(KeyMapContainer)+"/")] = v
		}
	}

	go func() {
		for {
			select {
			case x := <-ai.aChan:
				a.annotations[string(KeyMapContainer)+"/"+string(x.key)] = x.value
				ai.value[string(x.key)] = x.value
			case <-a.ctx.Done():
				return
			}
		}
	}()

	return ai
}

func (a MapContainer) Get(containerName string) (imageWithTag string, err error) {
	if v, ok := a.value[containerName]; ok {
		return v, nil
	}

	return "", fmt.Errorf("container %s not found", containerName)
}

// GetWithParser returns the custom image parser of the specified container.
func (a MapContainer) GetWithParser(containerName string) (imageWithTag utils.ImageTag, err error) {
	if v, ok := a.value[containerName]; ok {
		return utils.ImageParser(v), nil
	}

	return utils.ImageTag{}, fmt.Errorf("container %s not found", containerName)
}

func (a MapContainer) Set(containerName, imageWithTag string) {
	a.aChan.Send(AnnotationKey(containerName), imageWithTag)
}

func (a MapContainer) BuildPatches() (patches []patch.Patch) {
	for k, v := range a.value {
		patches = append(patches, patch.Patch{
			Op:    "add",
			Path:  "/metadata/annotations/" + strings.ReplaceAll(string(KeyMapContainer)+"."+k, "/", "~1"),
			Value: v,
		})
	}

	return
}

// * Generic funcs

func (aC aChan) Send(key AnnotationKey, value string) {
	aC <- struct {
		key   AnnotationKey
		value string
	}{key, value}
}
