-- name: CreateTransactions :one
INSERT INTO Transactions (
  to_account,
  from_account,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;
-- name: GetTransactions :one  
SELECT * FROM Transactions
WHERE id = $1 LIMIT 1;

-- name: ListTransactions :many
SELECT * FROM Transactions
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteTransactions :exec
DELETE FROM Transactions 
WHERE id=$1;