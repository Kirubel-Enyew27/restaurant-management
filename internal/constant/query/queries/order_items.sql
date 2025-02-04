-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, meal_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING order_item_id, order_id, meal_id, quantity, price;

-- name: GetOrderItemByID :one
SELECT order_item_id, order_id, meal_id, quantity, price
FROM order_items
WHERE order_item_id = $1;

-- name: GetOrderItemByOrderID :many
SELECT order_item_id, order_id, meal_id, quantity, price
FROM order_items
WHERE order_id = $1;

-- name: UpdateOrderItem :one
UPDATE order_items
SET 
    order_id = COALESCE(sqlc.narg('order_id'), order_id),
    meal_id = COALESCE(sqlc.narg('meal_id'), meal_id), 
    quantity = COALESCE(sqlc.narg('quantity'), quantity),
    price = COALESCE(sqlc.narg('price'), price)
WHERE order_item_id = $1
RETURNING order_item_id, order_id, meal_id, quantity, price;