package actions

import (
	"context"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
)

type (
	ActionInterface interface {
		Init(actualTag, newTag string, image *v1alpha1.Image)
		Execute(context.Context) error
	}

	Actions map[Name]ActionInterface
	Name    string

	action struct {
		actualTag string
		newTag    string
		image     *v1alpha1.Image
	}
)

var actions = make(Actions)

const (
	Apply Name = "apply"
)

func register(name Name, action ActionInterface) {
	actions[name] = action
}

// GetAction retrieves the ActionInterface associated with the given name.
// It takes a Name type as an argument and returns the corresponding ActionInterface.
// If the name does not exist in the actions map, the behavior is undefined.
//
// Parameters:
//   - name: The name of the action to retrieve.
//
// Returns://   - ActionInterface: The action associated with the provided name.n error indicating whether the action name was found or not. `ErrActionNotFound` is returned if the action name was not found.
func ParseActionName(name string) (Name, error) {
	for k := range actions {
		if k.String() == name {
			return Name(name), nil
		}
	}

	return "", ErrActionNotFound
}

func GetAction(name Name) (ActionInterface, error) {
	if _, ok := actions[name]; !ok {
		return nil, ErrActionNotFound
	}

	return actions[name], nil
}

// String returns the string representation of the action name.
func (n Name) String() string {
	return string(n)
}

func (a *action) Init(actualTag, newTag string, image *v1alpha1.Image) {
	a.actualTag = actualTag
	a.newTag = newTag
	a.image = image
}
