-- name: CreateMeal :one
INSERT INTO meals (name, price, img_url)
VALUES ($1, $2, $3)
RETURNING meal_id, name, img_url, price, available, created_at, modified_at;

-- name: GetMealByID :one
SELECT meal_id, name, img_url, price, available, created_at, modified_at
FROM meals
WHERE meal_id = $1;

-- name: GetMealByName :one
SELECT meal_id, name, img_url, price, available, created_at, modified_at
FROM meals
WHERE name = $1;

-- name: GetAllMeals :many
SELECT meal_id, name, img_url, price, available, created_at, modified_at
FROM meals;

-- name: UpdateMeal :one
UPDATE meals
SET
    name = COALESCE(sqlc.narg('name'), name),
    img_url = COALESCE(sqlc.narg('img_url'), img_url),
    price = COALESCE(sqlc.narg('price'), price),
    available = COALESCE(sqlc.narg('available'), available),
    modified_at = now()
WHERE meal_id = $1
RETURNING meal_id, name, img_url, price, available, created_at, modified_at;

-- name: DeleteMeal :exec
DELETE FROM meals
WHERE meal_id = $1;