package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/kube-image-updater/internal/rules"
)

func TestCalverMajor_Evaluate(t *testing.T) {
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
			actualTag:     "2024.01.00",
			tagsAvailable: []string{"2023.10.00", "2024.01.0", "2025.01.00"},
			expectedMatch: true,
			expectedTag:   "2025.01.00",
			expectedError: false,
		},
		{
			name:          "No matching major version",
			actualTag:     "2024.01.00",
			tagsAvailable: []string{"2023.10.00", "2024.10.01", "2024.01.10"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
		{
			name:          "Invalid actual tag",
			actualTag:     "invalid",
			tagsAvailable: []string{"2023.10.00", "2024.00.01", "2024.01.00"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: true,
		},
		{
			name:          "Invalid and no available tag",
			actualTag:     "2024.01.0",
			tagsAvailable: []string{"2023.01.0", "invalid"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := rules.GetRule(rules.CalverMajor)
			assert.NoError(t, err)
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

func TestCalverMinor_Evaluate(t *testing.T) {
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
			actualTag:     "2024.0.0",
			tagsAvailable: []string{"2024.1.0", "2024.0.1"},
			expectedMatch: true,
			expectedTag:   "2024.1.0",
			expectedError: false,
		},
		{
			name:          "No matching minor version",
			actualTag:     "2024.0.0",
			tagsAvailable: []string{"2024.0.1", "2024.0.2"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
		{
			name:          "Invalid actual tag",
			actualTag:     "invalid",
			tagsAvailable: []string{"2024.1.0", "2024.0.1"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: true,
		},
		{
			name:          "Invalid available tag",
			actualTag:     "2024.0.0",
			tagsAvailable: []string{"invalid", "2024.1.0"},
			expectedMatch: true,
			expectedTag:   "2024.1.0",
			expectedError: false,
		},
		{
			name:          "only major version increment",
			actualTag:     "2024.0.0",
			tagsAvailable: []string{"2025.0.0", "2024.0.0"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := rules.GetRule(rules.CalverMinor)
			assert.NoError(t, err)
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

func TestCalverPatch_Evaluate(t *testing.T) {
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
			actualTag:     "2024.0.0",
			tagsAvailable: []string{"2024.0.1", "2024.0.2"},
			expectedMatch: true,
			expectedTag:   "2024.0.2",
			expectedError: false,
		},
		{
			name:          "No matching patch version",
			actualTag:     "2024.0.0",
			tagsAvailable: []string{"2024.1.0", "2024.2.0"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
		{
			name:          "Invalid actual tag",
			actualTag:     "invalid",
			tagsAvailable: []string{"2024.0.1", "2024.0.2"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: true,
		},
		{
			name:          "Invalid available tag",
			actualTag:     "2024.0.0",
			tagsAvailable: []string{"invalid", "2024.0.1"},
			expectedMatch: true,
			expectedTag:   "2024.0.1",
			expectedError: false,
		},
		{
			name:          "only minor version increment",
			actualTag:     "2024.0.0",
			tagsAvailable: []string{"2024.1.0", "2024.0.0"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := rules.GetRule(rules.SemverPatch)
			assert.NoError(t, err)
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
