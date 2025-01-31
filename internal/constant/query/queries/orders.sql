-- name: CreateOrder :one
INSERT INTO orders (user_id, order_status, total_price)
VALUES ($1, $2, $3)
RETURNING order_id, user_id, order_status, total_price, created_at, modified_at;

-- name: GetOrderByID :one
SELECT order_id, user_id, order_status, total_price, created_at, modified_at
FROM orders
WHERE order_id = $1;

-- name: ListOrders :many
SELECT order_id, user_id, order_status, total_price, created_at, modified_at
FROM orders
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateOrderStatus :one
UPDATE orders
SET order_status = $2
WHERE order_id = $1
RETURNING order_id, user_id, order_status, total_price, created_at, modified_at;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE order_id = $1;
