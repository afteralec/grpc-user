-- name: CreateEmail :execresult
INSERT INTO emails (address, uid, verified) VALUES (?, ?, false);

-- name: MarkEmailVerified :exec
UPDATE emails SET verified = true WHERE id = ?;

-- name: GetEmail :one
SELECT * FROM emails WHERE id = ?;

-- name: GetEmailByAddressForUser :one
SELECT * FROM emails WHERE address = ? AND uid = ?;

-- name: GetVerifiedEmailByAddress :one
SELECT * FROM emails WHERE address = ? AND verified = true;

-- name: ListEmails :many
SELECT * FROM emails WHERE uid = ?;

-- name: ListVerifiedEmails :many
SELECT * FROM emails WHERE uid = ? AND verified = true;

-- name: CountEmails :one
SELECT COUNT(*) FROM emails WHERE uid = ?;

-- name: DeleteEmail :exec
DELETE FROM emails WHERE id = ?;
