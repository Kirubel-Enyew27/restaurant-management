-- name: CreateMeal :one
INSERT INTO meals (name, price, quantity)
VALUES ($1, $2, $3)
RETURNING meal_id, name, price, available, created_at, quantity, modified_at;

-- name: GetMealByID :one
SELECT meal_id, name, price, available, created_at, quantity, modified_at
FROM meals
WHERE meal_id = $1;

-- name: GetMealByName :one
SELECT meal_id, name, price, available, created_at, quantity, modified_at
FROM meals
WHERE name = $1;

-- name: GetAllMeals :many
SELECT meal_id, name, price, available, created_at, quantity, modified_at
FROM meals;

-- name: UpdateMeal :one
UPDATE meals
SET
    name = COALESCE(sqlc.narg('name'), name),
    price = COALESCE(sqlc.narg('price'), price),
    available = COALESCE(sqlc.narg('available'), available),
    quantity = COALESCE(sqlc.narg('quantity'), quantity),
    modified_at = now()
WHERE meal_id = $1
RETURNING meal_id, name, price, available, created_at, quantity, modified_at;

-- name: DeleteMeal :exec
DELETE FROM meals
WHERE meal_id = $1;