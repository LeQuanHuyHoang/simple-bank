-- name: CreateTransfer :one
INSERT INTO transfers (
    id,
    from_account_id,
    to_account_id,
    amount
) VALUES (
 $1, $2, $3, $4
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTransfer :many
SELECT * FROM transfers
WHERE from_account_id = $1 OR
      to_account_id = $2
ORDER BY id
LIMIT $3
OFFSET $4;

