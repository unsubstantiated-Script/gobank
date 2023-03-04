--Create
-- name: CreateAccount :one
INSERT INTO accounts (owner,
                      balance,
                      currency)
VALUES ($1, $2, $3)
RETURNING *;

--Read One
-- name: GetAccount :one
SELECT *
FROM accounts
WHERE id = $1
LIMIT 1;

-- Preventing concurrent queries from running at the same time NO KEY UPDATE keeps us out of concurrent locks
-- name: GetAccountForUpdate :one
SELECT *
FROM accounts
WHERE id = $1
LIMIT 1
FOR NO KEY UPDATE
;


--Read Many
-- name: ListAccounts :many
SELECT *
FROM accounts
ORDER BY id
LIMIT $1 OFFSET $2;

--Update
-- exec or one and RETURNING * to return an object if needed
-- name: UpdateAccount :one
UPDATE accounts
-- We're only updating the balance here. Doesn't make sense to update owner or currency...
SET balance  = $2
WHERE id = $1
RETURNING *;

--Add Balance
-- exec or one and RETURNING * to return an object if needed
-- name: AddAccountBalance :one
UPDATE accounts
-- We're only updating the balance here. Doesn't make sense to update owner or currency...
SET balance  = balance + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

--Delete
-- name: DeleteAccount :exec
DELETE FROM accounts WHERE id = $1;