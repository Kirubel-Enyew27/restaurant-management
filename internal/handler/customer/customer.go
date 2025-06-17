package customer

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"restaurant/internal/constant/errors"
	"restaurant/internal/constant/model/db"
	"restaurant/internal/constant/model/dto"
	"restaurant/internal/constant/model/response"
	"restaurant/internal/handler"
	"restaurant/internal/service"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type customer struct {
	logger         *zap.Logger
	customerModule service.Customer
	contextTimeout time.Duration
}

func Init(log *zap.Logger, customerModule service.Customer,
	contextTimeout time.Duration) handler.Customer {
	return &customer{
		logger:         log,
		customerModule: customerModule,
		contextTimeout: contextTimeout,
	}
}

func (cstmr *customer) Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	var req db.User

	if err := c.ShouldBindJSON(&req); err != nil {
		err := errors.ErrBadRequest.Wrap(err, "failed to bind request body")
		cstmr.logger.Info("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	newUser, err := cstmr.customerModule.Register(ctx, req)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusCreated, newUser, nil)
}

func (cstmr *customer) Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	var req db.User

	if err := c.ShouldBindJSON(&req); err != nil {
		err := errors.ErrBadRequest.Wrap(err, "failed to bind request body")
		cstmr.logger.Info("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	token, err := cstmr.customerModule.Login(ctx, db.User{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, token, nil)

}

func (cstmr *customer) GetCustomers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	users, err := cstmr.customerModule.GetUsers(ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, users, nil)

}

func (cstmr *customer) GetCustomerByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	userID := c.Param("id")
	user, err := cstmr.customerModule.GetUserByID(ctx, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, user, nil)

}

func (cstmr *customer) UpdateCustomer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	var reqBody dto.UpdateRequest

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		err := errors.ErrBadRequest.Wrap(err, "failed to bind request body")
		cstmr.logger.Info("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	userID := c.Param("id")

	updatedUser, err := cstmr.customerModule.UpdateUser(ctx, userID, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, updatedUser, nil)

}

func (cstmr *customer) DeleteCustomer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	userID := c.Param("id")

	err := cstmr.customerModule.DeleteUser(ctx, userID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, err, nil)

}

func (cstmr *customer) UploadProfilePicture(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	userID := c.Param("id")

	// Parse multipart form file
	file, header, err := c.Request.FormFile("profile_picture")
	if err != nil {
		err := errors.ErrBadRequest.Wrap(err, "failed to bind request body")
		cstmr.logger.Error("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}
	defer file.Close()

	// Simple file extension check
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		err := errors.ErrBadRequest.New("Only JPG, JPEG or PNG images are allowed")
		cstmr.logger.Error("invalid file extension", zap.Error(err))
		_ = c.Error(err)
		return
	}

	// Initialize Cloudinary client
	cloudName := os.Getenv("CLOUD_NAME")
	cloudApiKey := os.Getenv("CLOUD_API_KEY")
	cloudSecretKey := os.Getenv("CLOUD_SECRET_KEY")
	cld, err := cloudinary.NewFromParams(cloudName, cloudApiKey, cloudSecretKey)
	if err != nil {
		cstmr.logger.Error("Cloudinary init failed", zap.Error(err))
		_ = c.Error(errors.ErrInternalServerError.Wrap(err, "Cloudinary init failed"))
		return
	}

	// Upload file to Cloudinary
	overwrite := true
	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:  fmt.Sprint(userID),
		Folder:    "profile_pictures",
		Overwrite: &overwrite,
	})
	if err != nil {
		cstmr.logger.Error("Cloudinary upload failed", zap.Error(err))
		_ = c.Error(errors.ErrInternalServerError.Wrap(err, "Cloudinary upload failed"))
		return
	}

	imageURL := uploadResult.SecureURL

	// Update profile picture URL in DB
	reqBody := dto.UpdateRequest{ProfilePicture: imageURL}
	updatedUser, err := cstmr.customerModule.UpdateUser(ctx, userID, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, updatedUser, nil)
}

func (cstmr *customer) ChangePassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	var reqBody dto.ChangePassword

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		err := errors.ErrBadRequest.Wrap(err, "failed to bind request body")
		cstmr.logger.Info("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	userID := c.Param("id")

	updatedPassword, err := cstmr.customerModule.ChangePassword(ctx, userID, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, updatedPassword, nil)

}

func (cstmr *customer) SearchCustomer(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), cstmr.contextTimeout)
	defer cancel()

	query := c.Query("q")
	if query == "" {
		err := errors.ErrBadRequest.New("query param 'q' is required")
		cstmr.logger.Info("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	users, err := cstmr.customerModule.SearchUser(ctx, query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, users, nil)

}
