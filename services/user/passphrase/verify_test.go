package passphrase

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVerify(t *testing.T) {
	passphrase := "T3sted_tested"

	hash, err := Hash(passphrase, NewParams())
	if err != nil {
		t.Fatalf("Password hashing failed: %v", err)
	}

	verified, err := Verify(passphrase, hash)
	if err != nil {
		t.Fatalf("Password verification failed: %v", err)
	}

	require.True(t, verified)
}
