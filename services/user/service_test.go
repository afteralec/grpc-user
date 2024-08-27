package user

import (
	"context"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/afteralec/grpc-user/db"
)

const (
	TestUsername     = "testify"
	TestRootUsername = "tested"
	TestPassword     = "T3sted_tested"
)

func TestSyncRootPermissions(t *testing.T) {
	db, err := db.Open("../../test.db")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users;")
		db.Exec("DELETE FROM user_permission_grants;")
		db.Exec("DELETE FROM user_permission_revocations;")
		db.Close()
	})

	t.Run("CreatesRootUser", func(t *testing.T) {
		t.Cleanup(func() {
			db.Exec("DELETE FROM users;")
			db.Exec("DELETE FROM user_permission_grants;")
			db.Exec("DELETE FROM user_permission_revocations;")
		})
		config := viper.New()
		config.Set("root_username", TestRootUsername)
		config.Set("root_passphrase", TestPassword)
		ps, err := New(db, WithConfig(config))
		require.NoError(t, err)
		err = ps.SyncRootPermissions(context.Background())
		require.NoError(t, err)

		uid, err := ps.Authenticate(TestRootUsername, TestPassword)
		require.NoError(t, err)
		require.Greater(t, uid, int64(0))
	})

	t.Run("RevokesPermissionsFromPreviousRootUserAndGrantsToCurrent", func(t *testing.T) {
		t.Cleanup(func() {
			db.Exec("DELETE FROM users;")
			db.Exec("DELETE FROM user_permission_grants;")
			db.Exec("DELETE FROM user_permission_revocations;")
		})
		config := viper.New()
		config.Set("root_username", TestRootUsername)
		config.Set("root_passphrase", TestPassword)
		ps, err := New(db, WithConfig(config))
		require.NoError(t, err)
		ps.SyncRootPermissions(context.Background())
		uid, err := ps.Authenticate(TestRootUsername, TestPassword)
		require.NoError(t, err)
		require.Greater(t, uid, int64(0))

		config.Set("root_username", TestUsername)
		ps, err = New(db, WithConfig(config))
		require.NoError(t, err)
		err = ps.SyncRootPermissions(context.Background())
		require.NoError(t, err)

		records, err := ps.UserPermissions(context.Background(), uid)
		require.NoError(t, err)
		require.Empty(t, records)

		uid, err = ps.Authenticate(TestUsername, TestPassword)
		require.NoError(t, err)
		require.Greater(t, uid, int64(0))
		records, err = ps.UserPermissions(context.Background(), uid)
		require.NoError(t, err)
		require.NotEmpty(t, records)
		permissions := NewPermissions(uid, records)
		for _, permission := range RootPermissions {
			require.True(t, permissions.Has(permission.Name))
		}
	})
}

func TestRegister(t *testing.T) {
	db, err := db.Open("../../test.db")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users;")
		db.Exec("DELETE FROM user_permission_grants;")
		db.Exec("DELETE FROM user_permission_revocations;")
		db.Close()
	})

	config := viper.New()
	config.Set("root_username", TestRootUsername)
	ps, err := New(db, WithConfig(config))
	require.NoError(t, err)

	uid, err := ps.Register("testify", "T3sted_tested")
	require.NoError(t, err)
	require.NotEqual(t, 0, uid)
}

func TestRegisterAsRoot(t *testing.T) {
	db, err := db.Open("../../test.db")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users;")
		db.Exec("DELETE FROM user_permission_grants;")
		db.Exec("DELETE FROM user_permission_revocations;")
		db.Close()
	})

	config := viper.New()
	config.Set("root_username", TestRootUsername)

	ps, err := New(db, WithConfig(config))
	require.NoError(t, err)

	uid, err := ps.Register(TestRootUsername, "T3sted_tested")
	require.NoError(t, err)
	require.NotEqual(t, 0, uid)

	records, err := ps.UserPermissions(context.Background(), uid)
	require.NoError(t, err)
	permissions := NewPermissions(uid, records)
	for _, permission := range RootPermissions {
		require.True(t, permissions.Has(permission.Name))
	}
}

func TestAuthenticate(t *testing.T) {
	db, err := db.Open("../../test.db")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users;")
		db.Exec("DELETE FROM user_permission_grants;")
		db.Exec("DELETE FROM user_permission_revocations;")
		db.Close()
	})

	config := viper.New()
	config.Set("root_username", TestRootUsername)
	ps, err := New(db, WithConfig(config))
	require.NoError(t, err)

	uid, err := ps.Register("testify", "T3sted_tested")
	require.NoError(t, err)

	actual, err := ps.Authenticate("testify", "T3sted_tested")
	require.NoError(t, err)
	require.NotEqual(t, 0, uid)

	require.Equal(t, uid, actual)
}

func TestUsers(t *testing.T) {
	db, err := db.Open("../../test.db")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users;")
		db.Exec("DELETE FROM user_permission_grants;")
		db.Exec("DELETE FROM user_permission_revocations;")
		db.Close()
	})

	config := viper.New()
	config.Set("root_username", TestRootUsername)
	config.Set("root_passphrase", TestPassword)
	ps, err := New(db, WithConfig(config))
	require.NoError(t, err)
	err = ps.SyncRootPermissions(context.Background())
	require.NoError(t, err)

	users, err := ps.Users(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, 1, len(users))

	_, err = ps.Register(TestUsername, TestPassword)
	require.NoError(t, err)

	users, err = ps.Users(context.Background())
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.Equal(t, 2, len(users))
}

func TestUserSettings(t *testing.T) {
	db, err := db.Open("../../test.db")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users;")
		db.Exec("DELETE FROM user_permission_grants;")
		db.Exec("DELETE FROM user_permission_revocations;")
		db.Close()
	})

	config := viper.New()
	config.Set("root_username", TestRootUsername)
	ps, err := New(db, WithConfig(config))
	require.NoError(t, err)

	uid, err := ps.Register("testify", "T3sted_tested")
	require.NoError(t, err)

	settings, err := ps.UserSettings(context.Background(), uid)
	require.NoError(t, err)
	require.Equal(t, ThemeDefault, settings.Theme)
}

func TestSetUserSettingsTheme(t *testing.T) {
	db, err := db.Open("../../test.db")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users;")
		db.Exec("DELETE FROM user_permission_grants;")
		db.Exec("DELETE FROM user_permission_revocations;")
		db.Close()
	})

	config := viper.New()
	config.Set("root_username", TestRootUsername)
	ps, err := New(db, WithConfig(config))
	require.NoError(t, err)

	uid, err := ps.Register(TestUsername, TestPassword)
	require.NoError(t, err)

	settings, err := ps.UserSettings(context.Background(), uid)
	require.NoError(t, err)
	require.Equal(t, ThemeDefault, settings.Theme)

	settings, err = ps.SetUserSettingsTheme(context.Background(), uid, ThemeDark)
	require.NoError(t, err)
	require.Equal(t, ThemeDark, settings.Theme)

	settings, err = ps.UserSettings(context.Background(), uid)
	require.NoError(t, err)
	require.Equal(t, ThemeDark, settings.Theme)
}

func TestGrantAndRevokeUserPermissions(t *testing.T) {
	db, err := db.Open("../../test.db")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users;")
		db.Exec("DELETE FROM user_permission_grants;")
		db.Exec("DELETE FROM user_permission_revocations;")
		db.Close()
	})

	config := viper.New()
	config.Set("root_username", TestRootUsername)
	ps, err := New(db, WithConfig(config))
	require.NoError(t, err)

	iuid, err := ps.Register(TestRootUsername, TestPassword)
	require.NoError(t, err)
	uid, err := ps.Register(TestUsername, TestPassword)
	require.NoError(t, err)
	name := PermissionViewAllRooms.Name

	grantedID, err := ps.GrantUserPermission(context.Background(), uid, iuid, name)
	require.NoError(t, err)
	require.Greater(t, grantedID, int64(0))

	permissionRecords, err := ps.UserPermissions(context.Background(), uid)
	require.NoError(t, err)
	require.NotEmpty(t, permissionRecords)
	permissions := NewPermissions(uid, permissionRecords)
	require.True(t, permissions.Has(name))

	revokedID, err := ps.RevokeUserPermission(context.Background(), uid, iuid, name)
	require.NoError(t, err)
	require.Greater(t, revokedID, int64(0))
	require.Equal(t, grantedID, revokedID)

	permissionRecords, err = ps.UserPermissions(context.Background(), uid)
	require.NoError(t, err)
	require.Empty(t, permissionRecords)
	permissions = NewPermissions(uid, permissionRecords)
	require.False(t, permissions.Has(name))
}

func TestGrantAndRevokeUserPermissionReturnsEmptyIDForConflict(t *testing.T) {
	db, err := db.Open("../../test.db")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users;")
		db.Exec("DELETE FROM user_permission_grants;")
		db.Exec("DELETE FROM user_permission_revocations;")
		db.Close()
	})

	config := viper.New()
	config.Set("root_username", TestRootUsername)
	ps, err := New(db, WithConfig(config))
	require.NoError(t, err)

	iuid, err := ps.Register(TestRootUsername, TestPassword)
	require.NoError(t, err)
	uid, err := ps.Register(TestUsername, TestPassword)
	require.NoError(t, err)
	name := PermissionViewAllRooms.Name

	grantedID, err := ps.GrantUserPermission(context.Background(), uid, iuid, name)
	require.NoError(t, err)
	require.Greater(t, grantedID, int64(0))

	permissionRecords, err := ps.UserPermissions(context.Background(), uid)
	require.NoError(t, err)
	require.NotEmpty(t, permissionRecords)
	permissions := NewPermissions(uid, permissionRecords)
	require.True(t, permissions.Has(name))

	nextGrantedID, err := ps.GrantUserPermission(context.Background(), uid, iuid, name)
	require.NoError(t, err)
	require.Equal(t, int64(0), nextGrantedID)
}

func TestRevokeUserPermissionReturnsEmptyIDForConflict(t *testing.T) {
	db, err := db.Open("../../test.db")
	require.NoError(t, err)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users;")
		db.Exec("DELETE FROM user_permission_grants;")
		db.Exec("DELETE FROM user_permission_revocations;")
		db.Close()
	})

	config := viper.New()
	config.Set("root_username", TestRootUsername)
	ps, err := New(db, WithConfig(config))
	require.NoError(t, err)

	iuid, err := ps.Register(TestRootUsername, TestPassword)
	require.NoError(t, err)
	uid, err := ps.Register(TestUsername, TestPassword)
	require.NoError(t, err)
	name := PermissionViewAllRooms.Name

	grantedID, err := ps.GrantUserPermission(context.Background(), uid, iuid, name)
	require.NoError(t, err)
	require.Greater(t, grantedID, int64(0))

	permissionRecords, err := ps.UserPermissions(context.Background(), uid)
	require.NoError(t, err)
	require.NotEmpty(t, permissionRecords)
	permissions := NewPermissions(uid, permissionRecords)
	require.True(t, permissions.Has(name))

	nextGrantedID, err := ps.GrantUserPermission(context.Background(), uid, iuid, name)
	require.NoError(t, err)
	require.Equal(t, int64(0), nextGrantedID)

	revokedID, err := ps.RevokeUserPermission(context.Background(), uid, iuid, name)
	require.NoError(t, err)
	require.Greater(t, revokedID, int64(0))
	require.Equal(t, grantedID, revokedID)

	permissionRecords, err = ps.UserPermissions(context.Background(), uid)
	require.NoError(t, err)
	require.Empty(t, permissionRecords)
	permissions = NewPermissions(uid, permissionRecords)
	require.False(t, permissions.Has(name))

	nextRevokedID, err := ps.RevokeUserPermission(context.Background(), uid, iuid, name)
	require.NoError(t, err)
	require.Equal(t, int64(0), nextRevokedID)
}
