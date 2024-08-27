package user

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/afteralec/grpc-user/db/query"
)

func TestAllHardCodedPermissionsHaveValidNames(t *testing.T) {
	for _, permission := range AllPermissions {
		require.True(t, IsValidPermissionName(permission.Name))
	}
}

func TestIsValidNameReturnsFalseForInvalidName(t *testing.T) {
	require.False(t, IsValidPermissionName("not-a-permission"))
}

func TestNewPermissionsBuildsCorrectPermissions(t *testing.T) {
	uid := int64(1)
	permissionRecords := []query.UserPermission{
		{UID: uid, Name: PermissionGrantAll.Name},
		{UID: uid, Name: PermissionRevokeAll.Name},
	}
	permissions := NewPermissions(uid, permissionRecords)

	require.True(t, permissions.Has(PermissionGrantAll.Name))
	require.True(t, permissions.Has(PermissionRevokeAll.Name))

	for _, permission := range NonRootPermissions {
		require.False(t, permissions.Has(permission.Name))
	}
}

func TestCanGrant(t *testing.T) {
	uid := int64(1)
	records := []query.UserPermission{
		{UID: uid, Name: PermissionGrantAll.Name},
	}
	permissions := NewPermissions(uid, records)

	for _, permission := range NonRootPermissions {
		t.Run(fmt.Sprintf("can grant %s", permission.Name), func(t *testing.T) {
			require.True(t, permissions.CanGrant(permission.Name))
		})
	}
}

func TestCanGrantFalseForRootPermissions(t *testing.T) {
	uid := int64(1)
	records := []query.UserPermission{
		{UID: uid, Name: PermissionGrantAll.Name},
		{UID: uid, Name: PermissionRevokeAll.Name},
	}
	permissions := NewPermissions(uid, records)

	for _, permission := range RootPermissions {
		t.Run(fmt.Sprintf("cannot grant %s", permission.Name), func(t *testing.T) {
			require.False(t, permissions.CanGrant(permission.Name))
		})
	}
}

func TestCanRevokeFalseForRootPermissions(t *testing.T) {
	uid := int64(1)
	records := []query.UserPermission{
		{UID: uid, Name: PermissionGrantAll.Name},
		{UID: uid, Name: PermissionRevokeAll.Name},
	}
	permissions := NewPermissions(uid, records)

	for _, permission := range RootPermissions {
		t.Run(fmt.Sprintf("cannot revoke %s", permission.Name), func(t *testing.T) {
			require.False(t, permissions.CanRevoke(permission.Name))
		})
	}
}
