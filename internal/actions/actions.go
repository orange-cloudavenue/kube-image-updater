package actions

import (
	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/kubeclient"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
)

type (
	_actions map[models.ActionName]models.ActionInterface

	action struct {
		tags  models.Tags
		image *v1alpha1.Image
		k     kubeclient.Interface
		data  v1alpha1.ValueOrValueFrom
	}
)

var actions = make(_actions)

const (
	Apply        models.ActionName = "apply"
	AlertDiscord models.ActionName = "alert-discord"
	AlertEmail   models.ActionName = "alert-email"
)

func register(name models.ActionName, action models.ActionInterface) {
	actions[name] = action
}

// GetAction retrieves the ActionInterface associated with the given name.
// It takes a models.ActionName type as an argument and returns the corresponding ActionInterface.
// If the name does not exist in the actions map, the behavior is undefined.
//
// Parameters:
//   - name: The name of the action to retrieve.
//
// Returns://   - ActionInterface: The action associated with the provided name.n error indicating whether the action name was found or not. `ErrActionNotFound` is returned if the action name was not found.
func ParseActionName(name string) (models.ActionName, error) {
	for k := range actions {
		if k.String() == name {
			return models.ActionName(name), nil
		}
	}

	return "", ErrActionNotFound
}

// GetAction retrieves an action by its name.
// It returns the corresponding ActionInterface and an error if the action is not found.
//
// Parameters:
//   - name: The name of the action to retrieve.
//
// Returns:
//   - ActionInterface: The action associated with the given name.
//   - error: An error indicating if the action was not found (ErrActionNotFound).
func GetAction(name models.ActionName) (models.ActionInterface, error) {
	if _, ok := actions[name]; !ok {
		return nil, ErrActionNotFound
	}

	return actions[name], nil
}

// GetActionWithUntypedName retrieves an action based on the provided untyped name.
// It parses the action name and returns the corresponding ActionInterface.
// If the name cannot be parsed, it returns an error.
//
// Parameters:
//   - name: A string representing the untyped name of the action.
//
// Returns:
//   - An ActionInterface corresponding to the parsed action name, or nil if not found.
//   - An error if the action name could not be parsed.
func GetActionWithUntypedName(name string) (models.ActionInterface, error) {
	n, err := ParseActionName(name)
	if err != nil {
		return nil, err
	}
	return GetAction(n)
}

// * Generic action implementation

func (a *action) Init(kubeClient kubeclient.Interface, tags models.Tags, image *v1alpha1.Image, data v1alpha1.ValueOrValueFrom) {
	a.tags = tags
	a.image = image
	a.k = kubeClient
	a.data = data
}

// GetActualTag returns the current actual tag from the action's tags.
// It retrieves the value of the Actual field from the tags associated with the action.
func (a *action) GetActualTag() string {
	return a.tags.Actual
}

// GetNewTag returns the new tag from the action's tags.
// It retrieves the value of the New field from the tags associated with the action.
func (a *action) GetNewTag() string {
	return a.tags.New
}

// GetAvailableTags returns the available tags from the action's tags.
// It retrieves the value of the Available field from the tags associated with the action.
func (a *action) GetAvailableTags() []string {
	return a.tags.AvailableTags
}
