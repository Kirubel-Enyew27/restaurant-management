package food

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"restaurant/internal/constant/errors"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/constant/model/dto"
	"restaurant/internal/service"
	"restaurant/internal/storage"
	"strings"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type Food struct {
	log          *zap.Logger
	storage      storage.Food
	cacheStorage storage.OrderCache
	priceStorage storage.Price
}

func InitModule(
	log *zap.Logger,
	storage storage.Food,
	cache storage.FoodCache,
	priceStorage storage.Price,
) service.Food {
	return &Food{
		log:          log,
		storage:      storage,
		cacheStorage: cache,
		priceStorage: priceStorage,
	}
}

func (fd *Food) AddFood(ctx context.Context, data dto.FormData) (db.Meal, error) {
	if err := validation.ValidateStruct(&data,
		validation.Field(&data.Name, validation.Required),
		validation.Field(&data.Price, validation.Required),
		validation.Field(&data.File, validation.Required),
	); err != nil {
		fd.log.Error("failed to validate input", zap.Error(err))
		return db.Meal{}, errors.ErrInvalidUserInput.Wrap(err, "validation failed")
	}

	existingMeal, err := fd.storage.GetFoodByName(ctx, data.Name)
	if err != nil && !strings.Contains(err.Error(), pgx.ErrNoRows.Error()) {
		return db.Meal{}, err
	} else if existingMeal.Name != "" {
		fd.log.Error("food already exists", zap.Error(err))
		return db.Meal{}, errors.ErrDataAlredyExist.Wrap(err, "food already exists")
	}

	// Simple file extension check
	ext := strings.ToLower(filepath.Ext(data.Header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		err := errors.ErrBadRequest.New("Only JPG, JPEG or PNG images are allowed")
		fd.log.Error("invalid file extension", zap.Error(err))
		return db.Meal{}, err
	}

	// Initialize Cloudinary client
	cloudName := os.Getenv("CLOUD_NAME")
	cloudApiKey := os.Getenv("CLOUD_API_KEY")
	cloudSecretKey := os.Getenv("CLOUD_SECRET_KEY")
	cld, err := cloudinary.NewFromParams(cloudName, cloudApiKey, cloudSecretKey)
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "Cloudinary init failed")
		fd.log.Error("Cloudinary init failed", zap.Error(err))
		return db.Meal{}, err
	}

	// Upload file to Cloudinary
	overwrite := true
	uploadResult, err := cld.Upload.Upload(ctx, data.File, uploader.UploadParams{
		PublicID:  fmt.Sprint(uuid.New()),
		Folder:    "food_pictures",
		Overwrite: &overwrite,
	})
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "Cloudinary upload failed")
		fd.log.Error("Cloudinary upload failed", zap.Error(err))
		return db.Meal{}, err
	}

	imageURL := uploadResult.SecureURL

	meal := db.Meal{
		Name:   data.Name,
		Price:  data.Price,
		ImgUrl: imageURL,
	}

	return fd.storage.AddFood(ctx, meal)

}

func (fd *Food) GetFoods(ctx context.Context) ([]db.Meal, error) {
	return fd.storage.GetFoods(ctx)
}

func (fd *Food) GetFoodByID(ctx context.Context, foodID string) (db.Meal, error) {
	foodUUID, err := uuid.Parse(foodID)
	if err != nil {
		fd.log.Error("failed to parse food id", zap.Error(err))
		return db.Meal{}, errors.ErrInvalidUserInput.Wrap(err, "invalid food id")
	}

	return fd.storage.GetFoodByID(ctx, foodUUID)

}

func (fd *Food) UpdateFood(ctx context.Context, mealID string, req dto.FoodUpdate) (db.Meal, error) {
	mealUUID, err := uuid.Parse(mealID)
	if err != nil {
		fd.log.Error("failed to parse meal id", zap.Error(err))
		return db.Meal{}, errors.ErrInvalidUserInput.Wrap(err, "invalid meal id")
	}

	meal := db.Meal{
		MealID:    mealUUID,
		Name:      req.Name,
		Price:     req.Price,
		ImgUrl:    req.ImgUrl,
		Available: req.Available,
	}

	return fd.storage.UpdateFood(ctx, meal)
}

func (fd *Food) DeleteFood(ctx context.Context, mealID string) error {
	mealUUID, err := uuid.Parse(mealID)
	if err != nil {
		fd.log.Error("failed to parse meal id", zap.Error(err))
		return errors.ErrInvalidUserInput.Wrap(err, "invalid meal id")
	}

	meal, err := fd.storage.GetFoodByID(ctx, mealUUID)
	if err != nil {
		return err
	}

	return fd.storage.DeleteFood(ctx, meal.MealID)
}

func (fd *Food) SearchFood(ctx context.Context, query string) ([]db.Meal, error) {
	return fd.storage.SearchFood(ctx, query)
}
