package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/kube-image-updater/internal/rules"
)

func TestRegex_Evaluate(t *testing.T) {
	tests := []struct {
		name          string
		value         string
		actualTag     string
		tags          []string
		expectedMatch bool
		expectedTag   string
		expectError   bool
	}{
		{
			name:          "Match found",
			actualTag:     "v1.0.0",
			value:         `^v1\.0\.\d+$`,
			tags:          []string{"v1.0.1", "v1.0.2", "v1.1.0"},
			expectedMatch: true,
			expectedTag:   "v1.0.2",
			expectError:   false,
		},
		{
			name:          "No match found",
			actualTag:     "v1.0.0",
			value:         `^v2\.0\.\d+$`,
			tags:          []string{"v1.0.1", "v1.0.2", "v1.1.0"},
			expectedMatch: false,
			expectedTag:   "",
			expectError:   false,
		},
		{
			name:          "Invalid regex",
			actualTag:     "v1.0.0",
			value:         `^(v1\.0\.\d+$`,
			tags:          []string{"v1.0.1", "v1.0.2", "v1.1.0"},
			expectedMatch: false,
			expectedTag:   "",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rules.GetRule(rules.Regex)
			r.Init(tt.actualTag, tt.tags, tt.value)

			match, tag, err := r.Evaluate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedMatch, match)
			assert.Equal(t, tt.expectedTag, tag)
		})
	}
}
