package food

import (
	"context"
	"restaurant/internal/constant/errors"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/constant/model/persistencedb"
	"restaurant/internal/storage"

	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type food struct {
	db  persistencedb.PersistenceDB
	log *zap.Logger
}

func Init(db persistencedb.PersistenceDB, log *zap.Logger) storage.Food {
	return &food{
		db:  db,
		log: log,
	}
}

func (fd *food) AddFood(ctx context.Context, meal db.Meal) (db.Meal, error) {
	arg := db.CreateMealParams{
		Name:  meal.Name,
		Price: meal.Price,
	}

	newMeal, err := fd.db.Queries.CreateMeal(ctx, arg)
	if err != nil {
		fd.log.Info("failed to create meal", zap.Error(err))
		return db.Meal{}, errors.ErrUnableTocreate.Wrap(err, "failed to create user")
	}

	registeredMeal := db.Meal{
		MealID:     newMeal.MealID,
		Name:       newMeal.Name,
		Price:      newMeal.Price,
		Available:  newMeal.Available,
		CreatedAt:  newMeal.CreatedAt,
		ModifiedAt: newMeal.ModifiedAt,
	}

	return registeredMeal, nil
}

func (fd *food) GetFoods(ctx context.Context) ([]db.Meal, error) {
	meals, err := fd.db.Queries.GetAllMeals(ctx)
	if err != nil {
		if err == pgx.ErrNoRows {
			fd.log.Error("failed to get meals", zap.Error(err))
			return nil, errors.ErrUnableToGet.Wrap(err, "meals not found")
		}
		fd.log.Error("failed to get meals", zap.Error(err))
		return nil, errors.ErrUnableToGet.Wrap(err, "failed to get meals")
	}

	fetchedMeals := make([]db.Meal, len(meals))

	for i, meal := range meals {
		fetchedMeals[i] = db.Meal{
			MealID:     meal.MealID,
			Name:       meal.Name,
			Price:      meal.Price,
			Available:  meal.Available,
			CreatedAt:  meal.CreatedAt,
			ModifiedAt: meal.ModifiedAt,
		}
	}

	return fetchedMeals, nil
}

func (fd *food) GetFoodByName(ctx context.Context, name string) (db.Meal, error) {
	food, err := fd.db.Queries.GetMealByName(ctx, name)
	if err != nil {
		if err == pgx.ErrNoRows {
			fd.log.Error("failed to get meal by name", zap.Error(err))
			return db.Meal{}, errors.ErrNoRecordFound.Wrap(err, "meal not found")
		}
		fd.log.Error("failed to get meal by name", zap.Error(err))
		return db.Meal{}, errors.ErrUnableToGet.Wrap(err, "failed to get meal")
	}

	existingFood := db.Meal{
		MealID:     food.MealID,
		Name:       food.Name,
		Price:      food.Price,
		Available:  food.Available,
		CreatedAt:  food.CreatedAt,
		ModifiedAt: food.ModifiedAt,
	}

	return existingFood, nil
}
