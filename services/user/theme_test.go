package user

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsTheme(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected bool
	}
	testCases := [5]testCase{
		{"light", ThemeLight, true},
		{"dark", ThemeDark, true},
		{"default", ThemeDefault, true},
		{"random", strings.Repeat("a", 17), false},
		{"close", fmt.Sprintf("%s1234!", ThemeDefault), false},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.expected, IsTheme(tc.input), tc.name)
		})
	}
}
