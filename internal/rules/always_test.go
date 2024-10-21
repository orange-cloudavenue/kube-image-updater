package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/kube-image-updater/internal/rules"
)

func TestAlways_Evaluate(t *testing.T) {
	tests := []struct {
		name          string
		tags          []string
		actualTag     string
		expectedMatch bool
		expectedTag   string
	}{
		{
			name:          "New tag available",
			tags:          []string{"1.0.0", "1.1.0", "1.2.0"},
			actualTag:     "1.1.0",
			expectedMatch: true,
			expectedTag:   "1.2.0",
		},
		{
			name:          "No new tag available",
			tags:          []string{"1.0.0", "1.1.0"},
			actualTag:     "1.1.0",
			expectedMatch: false,
			expectedTag:   "",
		},
		{
			name:          "Empty tags",
			tags:          []string{},
			actualTag:     "1.0.0",
			expectedMatch: false,
			expectedTag:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := rules.GetRule(rules.Always)
			assert.NoError(t, err)
			r.Init(tt.actualTag, tt.tags, "")

			// err return always nil
			match, tag, _ := r.Evaluate()

			assert.Equal(t, tt.expectedMatch, match)
			assert.Equal(t, tt.expectedTag, tag)
		})
	}
}
