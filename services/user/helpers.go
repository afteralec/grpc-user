package user

import (
	"context"
	"database/sql"

	"github.com/afteralec/grpc-user/db/query"
)

func userPermissions(ctx context.Context, qtx *query.Queries, uid int64) ([]query.UserPermission, error) {
	permissions, err := qtx.ListUserPermissions(ctx, uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return []query.UserPermission{}, nil
		}
		return []query.UserPermission{}, err
	}
	return permissions, nil
}

func grantRootUserPermissions(ctx context.Context, qtx *query.Queries, uid int64) error {
	for _, permission := range RootPermissions {
		_, err := grantUserPermission(ctx, qtx, uid, 0, permission.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func userPermissionsByName(ctx context.Context, qtx *query.Queries, name string) ([]query.UserPermission, error) {
	permissions, err := qtx.ListUserPermissionsByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return []query.UserPermission{}, nil
		}
		return []query.UserPermission{}, err
	}
	return permissions, nil
}

func revokeAllRootUserPermissions(ctx context.Context, qtx *query.Queries) error {
	for _, permission := range RootPermissions {
		permissions, err := userPermissionsByName(ctx, qtx, permission.Name)
		if err != nil {
			return err
		}
		for _, permission := range permissions {
			if err := qtx.CreateUserPermissionRevocation(ctx, query.CreateUserPermissionRevocationParams{
				UID:  permission.UID,
				IUID: 0,
				Name: permission.Name,
			}); err != nil {
				return err
			}
		}
		if err := qtx.DeleteUserPermissionsByName(ctx, permission.Name); err != nil {
			return err
		}
	}
	return nil
}

func revokeUserPermission(ctx context.Context, qtx *query.Queries, uid, iuid int64, name string) (int64, error) {
	permission, err := qtx.GetUserPermissionByName(ctx, query.GetUserPermissionByNameParams{
		UID:  uid,
		Name: name,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	if err := qtx.DeleteUserPermission(ctx, permission.ID); err != nil {
		return 0, err
	}
	if err := qtx.CreateUserPermissionRevocation(ctx, query.CreateUserPermissionRevocationParams{
		UID:  uid,
		IUID: iuid,
		Name: name,
	}); err != nil {
		return 0, err
	}
	return permission.ID, nil
}

func grantUserPermission(ctx context.Context, qtx *query.Queries, uid, iuid int64, name string) (int64, error) {
	perms, err := userPermissions(ctx, qtx, uid)
	if err != nil {
		return 0, err
	}
	permissions := NewPermissions(uid, perms)
	if permissions.Has(name) {
		// TODO: Log this as "attempted to grant user permission they already have"
		return 0, nil
	}

	result, err := qtx.CreateUserPermission(ctx, query.CreateUserPermissionParams{
		UID:  uid,
		IUID: iuid,
		Name: name,
	})
	if err != nil {
		return 0, err
	}
	if err = qtx.CreateUserPermissionGrant(ctx, query.CreateUserPermissionGrantParams{
		UID:  uid,
		IUID: iuid,
		Name: name,
	}); err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}
