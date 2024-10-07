package actions

import (
	"context"

	"github.com/orange-cloudavenue/kube-image-updater/internal/annotations"
	"github.com/orange-cloudavenue/kube-image-updater/internal/models"
)

var _ models.ActionInterface = &apply{}

type (
	// apply is an action that applies the new tag to the image
	apply struct {
		action
	}
)

func init() {
	register(Apply, &apply{})
}

// Execute applies the new image tag to the image status.
// It returns an error if the new tag is empty.
//
// Parameters:
//   - ctx: The context for the operation.
//
// Returns:
//   - error: An error indicating the result of the operation, or nil if successful. `ErrEmptyNewTag` is returned if the new tag is empty.
func (a *apply) Execute(ctx context.Context) error {
	if a.GetNewTag() == "" {
		return ErrEmptyNewTag
	}

	an := annotations.New(ctx, a.image)
	an.Tag().Set(a.GetNewTag())

	// update the image with the new tag
	a.image.SetStatusTag(a.GetNewTag())

	return nil
}

// GetName returns the name of the action.
func (a *apply) GetName() models.ActionName {
	return Apply
}
