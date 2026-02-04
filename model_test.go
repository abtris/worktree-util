package main

import (
	"testing"
)

func TestSubstringFilter(t *testing.T) {
	tests := []struct {
		name     string
		term     string
		targets  []string
		expected int // expected number of matches
	}{
		{
			name:     "simple substring match",
			term:     "bug",
			targets:  []string{"bugfix", "feature", "debug"},
			expected: 2, // bugfix and debug
		},
		{
			name:     "case insensitive match",
			term:     "BUG",
			targets:  []string{"bugfix", "feature", "debug"},
			expected: 2,
		},
		{
			name:     "no matches",
			term:     "xyz",
			targets:  []string{"bugfix", "feature", "debug"},
			expected: 0,
		},
		{
			name:     "dash in term",
			term:     "-bugg",
			targets:  []string{"origin/ppoluektov/KR-7865-better-customer-facing-error-messages-in-direct-avs", "feature-buggy", "main"},
			expected: 1, // only feature-buggy contains "-bugg"
		},
		{
			name:     "partial match",
			term:     "feat",
			targets:  []string{"feature", "bugfix", "feature-auth"},
			expected: 2, // feature and feature-auth
		},
		{
			name:     "empty term",
			term:     "",
			targets:  []string{"feature", "bugfix"},
			expected: 2, // empty string matches everything
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ranks := substringFilter(tt.term, tt.targets)
			if len(ranks) != tt.expected {
				t.Errorf("substringFilter(%q, %v) returned %d matches, expected %d",
					tt.term, tt.targets, len(ranks), tt.expected)
				for _, rank := range ranks {
					t.Logf("  Matched: %s (index %d)", tt.targets[rank.Index], rank.Index)
				}
			}

			// Verify that matched indexes are set correctly
			for _, rank := range ranks {
				if len(rank.MatchedIndexes) == 0 && tt.term != "" {
					t.Errorf("Rank for %q has no matched indexes", tt.targets[rank.Index])
				}
			}
		})
	}
}

