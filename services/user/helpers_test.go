package user

import (
	"context"
	"testing"

	"github.com/afteralec/grpc-user/db"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

func TestRevokeAllRootPermissions(t *testing.T) {
	db, err := db.Open("../../test.db")
	require.NoError(t, err)
	viper.Set("root_username", TestRootUsername)
	t.Cleanup(func() {
		db.Exec("DELETE FROM users;")
		db.Exec("DELETE FROM user_permission_grants;")
		db.Exec("DELETE FROM user_permission_revocations;")
		db.Close()
	})

	config := viper.New()
	config.Set("root_username", TestRootUsername)
	service, err := New(db, WithConfig(config))
	require.NoError(t, err)

	uid, err := service.Register(TestRootUsername, "T3sted_tested")
	require.NoError(t, err)
	require.NotEqual(t, 0, uid)

	records, err := service.UserPermissions(context.Background(), uid)
	require.NoError(t, err)
	permissions := NewPermissions(uid, records)
	for _, permission := range RootPermissions {
		require.True(t, permissions.Has(permission.Name))
	}

	tx, err := service.db.Begin()
	require.NoError(t, err)
	t.Cleanup(func() {
		tx.Rollback()
	})
	qtx := service.query.WithTx(tx)

	err = revokeAllRootUserPermissions(context.Background(), qtx)
	require.NoError(t, err)

	err = tx.Commit()
	require.NoError(t, err)

	records, err = service.UserPermissions(context.Background(), uid)
	require.NoError(t, err)
	permissions = NewPermissions(uid, records)
	for _, permission := range RootPermissions {
		require.False(t, permissions.Has(permission.Name))
	}
}
