-- name: CreateUser :exec
INSERT INTO users (username, password_hash)
VALUES ($1, $2);

-- name: GetUserByUsername :one
SELECT username, password_hash, coins
FROM users
WHERE username = $1;

-- name: DeductCoins :execrows
UPDATE users
SET coins = coins - $1
WHERE username = $2 AND coins >= $1;

-- name: AddCoins :execrows
UPDATE users
SET coins = coins + $1
WHERE username = $2;

-- name: InsertCoinTransfer :exec
INSERT INTO coin_transfers (from_username, to_username, amount)
VALUES ($1, $2, $3);

-- name: GetCoinHistoryReceived :many
SELECT from_username, amount, created_at
FROM coin_transfers
WHERE to_username = $1
ORDER BY created_at;

-- name: GetCoinHistorySent :many
SELECT to_username, amount
FROM coin_transfers
WHERE from_username = $1;

-- name: GetProductPrice :one
SELECT price
FROM products
WHERE item = $1;

-- name: CreatePurchase :one
INSERT INTO purchases (username, item, price)
VALUES ($1, $2, $3)
RETURNING id, username, item, price, created_at;

-- name: ListInventory :many
SELECT item, COUNT(*) AS quantity
FROM purchases
WHERE username = $1
GROUP BY item;

-- name: GetUser :one
SELECT username, password_hash, coins
FROM users
WHERE username = $1;