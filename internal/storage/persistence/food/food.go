package food

import (
	"context"
	"database/sql"
	"restaurant/internal/constant/errors"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/constant/model/persistencedb"
	"restaurant/internal/storage"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/shopspring/decimal"
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
		return db.Meal{}, errors.ErrUnableTocreate.Wrap(err, "failed to create meal")
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

func (fd *food) GetFoodByID(ctx context.Context, mealID uuid.UUID) (db.Meal, error) {
	meal, err := fd.db.Queries.GetMealByID(ctx, mealID)
	if err != nil {
		if err == pgx.ErrNoRows {
			fd.log.Error("failed to get meal by id", zap.Error(err))
			return db.Meal{}, errors.ErrNoRecordFound.Wrap(err, "meal not found")
		}
		fd.log.Error("failed to get meal by id", zap.Error(err))
		return db.Meal{}, errors.ErrUnableToGet.Wrap(err, "failed to get meal")
	}

	Meal := db.Meal{
		MealID:     meal.MealID,
		Name:       meal.Name,
		Price:      meal.Price,
		Available:  meal.Available,
		CreatedAt:  meal.CreatedAt,
		ModifiedAt: meal.ModifiedAt,
		ImgUrl:     meal.ImgUrl,
	}

	return Meal, nil

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

func (fd *food) UpdateFood(ctx context.Context, meal db.Meal) (db.Meal, error) {
	updateParams := db.UpdateMealParams{
		MealID:    meal.MealID,
		Name:      sql.NullString{},
		ImgUrl:    sql.NullString{},
		Price:     decimal.NullDecimal{},
		Available: sql.NullBool{},
	}

	// Set values if non-empty
	if meal.Name != "" {
		updateParams.Name = sql.NullString{String: meal.Name, Valid: true}
	}
	if meal.ImgUrl.String != "" {
		updateParams.ImgUrl = sql.NullString{String: meal.ImgUrl.String, Valid: true}
	}
	if meal.Price.IsZero() == false {
		updateParams.Price = decimal.NullDecimal{Decimal: meal.Price, Valid: true}
	}
	if meal.Available.Valid {
		updateParams.Available = sql.NullBool{Bool: meal.Available.Bool, Valid: true}
	}

	Meal, err := fd.db.Queries.UpdateMeal(ctx, updateParams)
	if err != nil {
		if err == pgx.ErrNoRows {
			fd.log.Error("meal to be updated does not exist", zap.Error(err))
			return db.Meal{}, errors.ErrNoRecordFound.Wrap(err, "meal not found")
		}
		fd.log.Error("failed to update meal", zap.Error(err))
		return db.Meal{}, errors.ErrUnableToUpdate.Wrap(err, "failed to update meal")
	}

	updatedMeal := db.Meal{
		MealID:     Meal.MealID,
		Name:       Meal.Name,
		Price:      Meal.Price,
		Available:  Meal.Available,
		CreatedAt:  Meal.CreatedAt,
		ModifiedAt: Meal.ModifiedAt,
		ImgUrl:     Meal.ImgUrl,
	}

	return updatedMeal, nil
}

func (fd *food) DeleteFood(ctx context.Context, mealID uuid.UUID) error {
	err := fd.db.Queries.DeleteMeal(ctx, mealID)
	if err != nil {
		fd.log.Error("error deleting meal", zap.Error(err))
		return errors.ErrDBDelError.Wrap(err, "failed to delete meal")
	}

	return nil
}
