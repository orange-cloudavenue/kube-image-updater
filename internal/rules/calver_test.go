package rules_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/orange-cloudavenue/kube-image-updater/internal/rules"
)

var listTest map[string]string = map[string]string{
	"YYYY":                       "2024",
	"YYYY-dev":                   "2024-dev",
	"YYYY.MM":                    "2024.01",
	"YYYY.MM-dev":                "2024.01-dev",
	"YYYY.MM.DD":                 "2024.01.01",
	"YYYY.MM.DD-dev":             "2024.01.01-dev",
	"YYYY.MM.DD-dev.prerelease":  "2024.01.01-dev.1",
	"YY":                         "24",
	"YY-dev":                     "24-dev",
	"YY.MM":                      "24.01",
	"YY.M":                       "24.1",
	"YY.M.D":                     "24.1.1",
	"Invalid":                    "v2024",
	"YYYY.MM.DD.WrongPrerelease": "2024.01.01.1",
}

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
			// Unitary tests
			name:          "YYYY",
			actualTag:     listTest["YYYY"],
			tagsAvailable: []string{"2023", "2024", "2025"},
			expectedMatch: true,
			expectedTag:   "2025",
			expectedError: false,
		},
		{
			name:          "YY",
			actualTag:     listTest["YY"],
			tagsAvailable: []string{"23", "24", "25"},
			expectedMatch: true,
			expectedTag:   "25",
			expectedError: false,
		},
		{
			name:          "YYYY-dev",
			actualTag:     listTest["YYYY-dev"],
			tagsAvailable: []string{"2023-dev", "2024-dev", "2025-dev"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
		{
			name:          "YY-dev",
			actualTag:     listTest["YY-dev"],
			tagsAvailable: []string{"23-dev", "24-dev", "25-dev"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
		// Errors tests
		{
			name:          "Invalid",
			actualTag:     listTest["Invalid"],
			tagsAvailable: []string{"2023.10.00", "2024.01.0", "2025.01.00"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: true,
		},
		{
			name:          "Invalid available tag",
			actualTag:     "2024.01.0",
			tagsAvailable: []string{"v2023.01.0", "invalid"},
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
			name:          "YYYY.MM",
			actualTag:     listTest["YYYY.MM"],
			tagsAvailable: []string{"2023.01", "2025.01", "2024.02"},
			expectedMatch: true,
			expectedTag:   "2024.02",
			expectedError: false,
		},
		{
			name:          "YYYY.MM-dev",
			actualTag:     listTest["YYYY.MM-dev"],
			tagsAvailable: []string{"2023.01-dev", "2025.01-dev", "2024.02-dev"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
		{
			name:          "YY.MM",
			actualTag:     listTest["YY.MM"],
			tagsAvailable: []string{"23.01", "25.01", "24.02"},
			expectedMatch: true,
			expectedTag:   "24.02",
			expectedError: false,
		},
		{
			name:          "YY.M",
			actualTag:     listTest["YY.M"],
			tagsAvailable: []string{"23.1", "25.1", "24.2"},
			expectedMatch: true,
			expectedTag:   "24.2",
			expectedError: false,
		},
		// Errors tests
		{
			name:          "Invalid",
			actualTag:     listTest["Invalid"],
			tagsAvailable: []string{"2023.10.00", "2024.01.0", "2025.01.00"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: true,
		},
		{
			name:          "Invalid available tag",
			actualTag:     "2024.01.0",
			tagsAvailable: []string{"v2023.01.0", "invalid"},
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
		// Unitary tests
		{
			name:          "YYYY.MM.DD",
			actualTag:     listTest["YYYY.MM.DD"],
			tagsAvailable: []string{"2023.01.01", "2025.02.01", "2024.01.02"},
			expectedMatch: true,
			expectedTag:   "2024.01.02",
			expectedError: false,
		},
		{
			name:          "YYYY.MM.DD-dev",
			actualTag:     listTest["YYYY.MM.DD-dev"],
			tagsAvailable: []string{"2023.01.01-dev", "2025.02.01-dev", "2024.01.02-dev"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
		{
			name:          "YY.M.D",
			actualTag:     listTest["YY.M.D"],
			tagsAvailable: []string{"23.1.1", "25.2.1", "24.1.2"},
			expectedMatch: true,
			expectedTag:   "24.1.2",
			expectedError: false,
		},
		// Errors tests
		{
			name:          "YY.M.D-dev",
			actualTag:     listTest["YY.M.D-dev"],
			tagsAvailable: []string{"23.1.1-dev", "25.2.1-dev", "24.1.2-dev"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: true,
		},
		{
			name:          "Invalid",
			actualTag:     listTest["Invalid"],
			tagsAvailable: []string{"2023.10.00", "2024.01.0", "2025.01.00"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: true,
		},
		{
			name:          "Invalid available tag",
			actualTag:     "2024.01.0",
			tagsAvailable: []string{"v2023.01.0", "invalid"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := rules.GetRule(rules.CalverPatch)
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

// func TestCalverModifier_Evaluate(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		actualTag     string
// 		tagsAvailable []string
// 		expectedMatch bool
// 		expectedTag   string
// 		expectedError bool
// 	}{
// 		// Unitary tests
// 		{
// 			name:          "YYYY-dev",
// 			actualTag:     listTest["YYYY-dev"],
// 			tagsAvailable: []string{"2024-aaa", "2024-dev.1", "2025-dev"},
// 			expectedMatch: true,
// 			expectedTag:   "2025-dev",
// 			expectedError: false,
// 		},
// 		{
// 			name:          "YYYY.MM-dev",
// 			actualTag:     listTest["YYYY.MM-dev"],
// 			tagsAvailable: []string{"2024.01-aaa", "2024.02-bbb", "2024.02-dev"},
// 			expectedMatch: true,
// 			expectedTag:   "2024.02-dev",
// 			expectedError: false,
// 		},
// 		{
// 			name:          "YYYY.MM.DD-dev",
// 			actualTag:     listTest["YYYY.MM.DD-dev"],
// 			tagsAvailable: []string{"2024.01.01-aaa", "2024.01.01-dev", "2024.01.02-dev"},
// 			expectedMatch: true,
// 			expectedTag:   "2024.01.02-dev",
// 			expectedError: false,
// 		},
// 		{
// 			name:          "YY-dev",
// 			actualTag:     listTest["YY-dev"],
// 			tagsAvailable: []string{"24-aaa", "24-dev", "25-dev"},
// 			expectedMatch: true,
// 			expectedTag:   "25-dev",
// 			expectedError: false,
// 		},
// 		{
// 			name:          "YY.MM-dev",
// 			actualTag:     listTest["YY.MM-dev"],
// 			tagsAvailable: []string{"24.01-aaa", "24.01-dev", "24.02-dev"},
// 			expectedMatch: true,
// 			expectedTag:   "24.02-dev",
// 			expectedError: false,
// 		},
// 		{
// 			name:          "YY.M.D-dev",
// 			actualTag:     listTest["YY.M.D-dev"],
// 			tagsAvailable: []string{"24.1.1-aaa", "24.1.1-dev", "24.1.2-dev"},
// 			expectedMatch: true,
// 			expectedTag:   "24.1.2-dev",
// 			expectedError: false,
// 		},
// 		// Errors tests
// 		{
// 			name:          "Invalid",
// 			actualTag:     listTest["Invalid"],
// 			tagsAvailable: []string{"2023.10.00-dev", "24.01.0-dev", "2025.1.0-dev"},
// 			expectedMatch: false,
// 			expectedTag:   "",
// 			expectedError: true,
// 		},
// 		{
// 			name:          "Invalid available tag",
// 			actualTag:     "2024.01.0",
// 			tagsAvailable: []string{"v2023.01.0", "invalid"},
// 			expectedMatch: false,
// 			expectedTag:   "",
// 			expectedError: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r, err := rules.GetRule(rules.CalverModifier)
// 			assert.NoError(t, err)
// 			r.Init(tt.actualTag, tt.tagsAvailable, "")
// 			match, newTag, err := r.Evaluate()

// 			if tt.expectedError {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 			}

// 			assert.Equal(t, tt.expectedMatch, match)
// 			assert.Equal(t, tt.expectedTag, newTag)
// 		})
// 	}
// }

func TestCalverPrerelease_Evaluate(t *testing.T) {
	tests := []struct {
		name          string
		actualTag     string
		tagsAvailable []string
		expectedMatch bool
		expectedTag   string
		expectedError bool
	}{
		// Unitary tests
		{
			name:          "YYYY.MM.DD-dev.prerelease",
			actualTag:     listTest["YYYY.MM.DD-dev.prerelease"],
			tagsAvailable: []string{"2023.01.01-dev.2", "2024.01.01-beta.2", "2024.01.01-dev.2"},
			expectedMatch: true,
			expectedTag:   "2024.01.01-dev.2",
			expectedError: false,
		},
		// Errors tests
		{
			name:          "YYYY.MM.DD.WrongPrerelease",
			actualTag:     listTest["YYYY.MM.DD.WrongPrerelease"],
			tagsAvailable: []string{"2024.01.01-aaa", "2024.01.01-dev", "2024.01.01.2"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: true,
		},
		{
			name:          "Invalid",
			actualTag:     listTest["Invalid"],
			tagsAvailable: []string{"2023.10.00-dev", "24.01.0-dev", "2025.1.0-dev.2"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: true,
		},
		{
			name:          "Invalid available tag",
			actualTag:     "2024.01.0",
			tagsAvailable: []string{"v2023.01.0.aaa", "invalid.prelease"},
			expectedMatch: false,
			expectedTag:   "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := rules.GetRule(rules.CalverPrerelease)
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
