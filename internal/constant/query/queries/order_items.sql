-- name: CreateOrderItem :one
INSERT INTO order_items (order_id, meal_id, quantity, price)
VALUES ($1, $2, $3, $4)
RETURNING order_item_id, order_id, meal_id, quantity, price;

-- name: GetOrderItemByID :one
SELECT order_item_id, order_id, meal_id, quantity, price
FROM order_items
WHERE order_item_id = $1;