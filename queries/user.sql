-- name: CreateUser :execresult
INSERT INTO users (username, pw_hash) VALUES (?, ?);

-- name: UpdateUserPassword :execresult
UPDATE users SET pw_hash = ? WHERE id = ?;

-- name: ListUsers :many
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users WHERE id = ?;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = ?;

-- name: GetUserUsername :one
SELECT username FROM users WHERE id = ?;

-- name: GetUSerUsernameById :one
SELECT username FROM users WHERE id = ?;

-- name: SearchUsersByUsername :many
SELECT * FROM users WHERE username LIKE ?;

-- name: CreateUserPermission :execresult
INSERT INTO user_permissions (name, uid, iuid) VALUES (?, ?, ?);

-- name: GetUserPermissionByName :one
SELECT * FROM user_permissions WHERE name = ? AND uid = ?;

-- name: DeleteUserPermission :exec
DELETE FROM user_permissions WHERE id = ?;

-- name: ListUserPermissionsByName :many
SELECT * FROM user_permissions WHERE name = ?;

-- name: DeleteUserPermissionsByName :exec
DELETE FROM user_permissions WHERE name = ?;

-- name: ListUserPermissions :many
SELECT * FROM user_permissions WHERE uid = ?;

-- name: CreateUserPermissionGrant :exec
INSERT INTO user_permission_grants (name, uid, iuid) VALUES (?, ?, ?);

-- name: CreateUserPermissionRevocation :exec
INSERT INTO user_permission_revocations (name, uid, iuid) VALUES (?, ?, ?);

-- name: CreateUserSettings :exec
INSERT INTO user_settings (theme, uid) VALUES (?, ?);

-- name: GetUserSettings :one
SELECT * FROM user_settings WHERE uid = ?;

-- name: UpdateUserSettingsTheme :exec
UPDATE user_settings SET theme = ? WHERE uid = ?;
