CREATE TABLE "transfers" (
    "id" bigserial PRIMARY KEY,
    "from_account_id" bigint NOT NULL,
    "to_account_id" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())   
);


-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount,
    created_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetTransfer :one
SELECT * 
FROM transfers
WHERE id=$1
LIMIT 1;

-- name: ListTransfer :many
SELECT *
FROM transfers
WHERE from_account_id = $1 OR to_account_id = $2
ORDER BY id
LIMIT $1
OFFSET $2; --смещение

-- name: UpdateTransfer :one
UPDATE transfers
SET amount=$1 AND to_account_id = $2
WHERE id=$2
RETURNING *;

-- name: DeleteTransfer :exec
DELETE FROM transfers
WHERE id = $1;