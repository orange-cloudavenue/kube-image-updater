package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/kube-image-updater/internal/rules"
)

func TestSemverMajor_Evaluate(t *testing.T) {
	tests := []struct {
		name          string
		actualTag     string
		tagsAvailable []string
		expectedMatch bool
		expectedTag   string
		expectedError bool
	}{
		{
			name:          "Valid major version increment",
			actualTag:     "1.0.0",
			tagsAvailable: []string{"2.0.0", "1.1.0"},
			expectedMatch: true,
			expectedTag:   "2.0.0",
			expectedError: false,
		},
		{
			name:          "No matching major version",
			actualTag:     "1.0.0",
			tagsAvailable: []string{"1.1.0", "1.2.0"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
		{
			name:          "Invalid actual tag",
			actualTag:     "invalid",
			tagsAvailable: []string{"2.0.0", "1.1.0"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: true,
		},
		{
			name:          "Invalid available tag",
			actualTag:     "1.0.0",
			tagsAvailable: []string{"invalid", "2.0.0"},
			expectedMatch: true,
			expectedTag:   "2.0.0",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rules.GetRule(rules.SemverMajor)
			r.Init(tt.actualTag, tt.tagsAvailable, "")
			match, newTag, err := r.Evaluate()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedMatch, match)
			assert.Equal(t, tt.expectedTag, newTag)
		})
	}
}

func TestSemverMinor_Evaluate(t *testing.T) {
	tests := []struct {
		name          string
		actualTag     string
		tagsAvailable []string
		expectedMatch bool
		expectedTag   string
		expectedError bool
	}{
		{
			name:          "Valid minor version increment",
			actualTag:     "1.0.0",
			tagsAvailable: []string{"1.1.0", "1.0.1"},
			expectedMatch: true,
			expectedTag:   "1.1.0",
			expectedError: false,
		},
		{
			name:          "No matching minor version",
			actualTag:     "1.0.0",
			tagsAvailable: []string{"1.0.1", "1.0.2"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
		{
			name:          "Invalid actual tag",
			actualTag:     "invalid",
			tagsAvailable: []string{"1.1.0", "1.0.1"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: true,
		},
		{
			name:          "Invalid available tag",
			actualTag:     "1.0.0",
			tagsAvailable: []string{"invalid", "1.1.0"},
			expectedMatch: true,
			expectedTag:   "1.1.0",
			expectedError: false,
		},
		{
			name:          "only major version increment",
			actualTag:     "1.0.0",
			tagsAvailable: []string{"2.0.0", "1.0.0"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rules.GetRule(rules.SemverMinor)
			r.Init(tt.actualTag, tt.tagsAvailable, "")
			match, newTag, err := r.Evaluate()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedMatch, match)
			assert.Equal(t, tt.expectedTag, newTag)
		})
	}
}

func TestSemverPatch_Evaluate(t *testing.T) {
	tests := []struct {
		name          string
		actualTag     string
		tagsAvailable []string
		expectedMatch bool
		expectedTag   string
		expectedError bool
	}{
		{
			name:          "Valid patch version increment",
			actualTag:     "1.0.0",
			tagsAvailable: []string{"1.0.1", "1.0.2"},
			expectedMatch: true,
			expectedTag:   "1.0.2",
			expectedError: false,
		},
		{
			name:          "No matching patch version",
			actualTag:     "1.0.0",
			tagsAvailable: []string{"1.1.0", "1.2.0"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
		{
			name:          "Invalid actual tag",
			actualTag:     "invalid",
			tagsAvailable: []string{"1.0.1", "1.0.2"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: true,
		},
		{
			name:          "Invalid available tag",
			actualTag:     "1.0.0",
			tagsAvailable: []string{"invalid", "1.0.1"},
			expectedMatch: true,
			expectedTag:   "1.0.1",
			expectedError: false,
		},
		{
			name:          "only minor version increment",
			actualTag:     "1.0.0",
			tagsAvailable: []string{"1.1.0", "1.0.0"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := rules.GetRule(rules.SemverPatch)
			r.Init(tt.actualTag, tt.tagsAvailable, "")
			match, newTag, err := r.Evaluate()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedMatch, match)
			assert.Equal(t, tt.expectedTag, newTag)
		})
	}
}
