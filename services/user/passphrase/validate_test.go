package passphrase

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValid(t *testing.T) {
	type testCase struct {
		name  string
		input string
		want  bool
	}
	testCases := [5]testCase{
		{"too short", strings.Repeat("a", 3), false},
		{"normal", "T3sted tested", true},
		{"special characters", "!@#$%^&*", true},
		{"long normal", fmt.Sprintf("%s%s", "T3sted tested", strings.Repeat("a", 242)), true},
		{"too long", strings.Repeat("a", 256), false},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			require.Equal(t, tc.want, IsValid(tc.input), tc.name)
		})
	}
}
