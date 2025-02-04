-- name: CreateOrder :one
INSERT INTO orders (user_id, order_status, total_price)
VALUES ($1, $2, $3)
RETURNING order_id, user_id, order_status, total_price, created_at, modified_at;

-- name: GetOrderByID :one
SELECT order_id, user_id, order_status, total_price, created_at, modified_at
FROM orders
WHERE order_id = $1;

-- name: ListOrders :many
	SELECT 
		o.order_id,
		o.user_id,
		o.order_status,
		o.total_price,
		o.created_at,
		o.modified_at,
		oi.order_item_id,
		oi.meal_id,
		oi.quantity,
		oi.price
	FROM orders o
	LEFT JOIN order_items oi ON o.order_id = oi.order_id
	ORDER BY o.created_at DESC;

-- name: UpdateOrder :one
UPDATE orders
SET 
    order_status = COALESCE(sqlc.narg('order_status'), order_status), 
    total_price = COALESCE(sqlc.narg('total_price'), total_price), 
    modified_at = now()
WHERE order_id = $1
RETURNING order_id, user_id, order_status, total_price, created_at, modified_at;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE order_id = $1;
