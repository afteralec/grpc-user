package passphrase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHash(t *testing.T) {
	passphrase := "T3sted_tested"

	hash, err := Hash(passphrase, NewParams())
	if err != nil {
		t.Fatalf("Password hashing failed: %v", err)
	}

	require.Greater(t, len(hash), 0)
}
