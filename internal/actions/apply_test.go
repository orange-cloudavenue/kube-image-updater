package actions_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/kube-image-updater/api/v1alpha1"
	"github.com/orange-cloudavenue/kube-image-updater/internal/actions"
)

func TestApply_Execute(t *testing.T) {
	tests := []struct {
		name        string
		initialTag  string
		newTag      string
		expectedTag string
		expectedErr error
	}{
		{
			name:        "Valid tag update",
			initialTag:  "1.0.0",
			newTag:      "1.1.0",
			expectedTag: "1.1.0",
			expectedErr: nil,
		},
		{
			name:        "Empty new tag",
			initialTag:  "1.0.0",
			newTag:      "",
			expectedTag: "",
			expectedErr: actions.ErrEmptyNewTag,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := actions.GetAction(actions.Apply)
			assert.NoError(t, err)
			image := &v1alpha1.Image{}

			a.Init(tt.initialTag, tt.newTag, image)

			err = a.Execute(context.Background())
			if tt.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedTag, image.Status.Tag)
			}
		})
	}
}
