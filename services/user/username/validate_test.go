package username

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValid(t *testing.T) {
	type testcase struct {
		name        string
		input       string
		expectError bool
	}
	testcases := [4]testcase{
		{"too short", strings.Repeat("a", 3), true},
		{"normal", "test", false},
		{"normal but long", fmt.Sprintf("%s%s", "test", strings.Repeat("a", 10)), false},
		{"too long", strings.Repeat("a", 17), true},
	}
	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := IsValid(tc.input)
			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
