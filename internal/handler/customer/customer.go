package customer

import (
	"context"
	"fmt"
	"io"
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
		cstmr.logger.Error("invalid request body", zap.Error(err))
		_ = c.Error(err)
		return
	}

	// Create destination path
	uploadDir := "uploads/"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}
	filename := fmt.Sprintf("%s%s%s", uploadDir, userID, ext)

	// Save file
	out, err := os.Create(filename)
	if err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "Unable to save file")
		cstmr.logger.Error("Unable to save file", zap.Error(err))
		_ = c.Error(err)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		err := errors.ErrInternalServerError.Wrap(err, "Failed to write file")
		cstmr.logger.Error("Failed to write file", zap.Error(err))
		_ = c.Error(err)
		return
	}

	// Update profile picture in the database
	imagePath := fmt.Sprintf("/%s", filename)
	reqBody := dto.UpdateRequest{ProfilePicture: imagePath}
	updatedUser, err := cstmr.customerModule.UpdateUser(ctx, userID, reqBody)
	if err != nil {
		_ = c.Error(err)
		return
	}

	response.SendSuccessResponse(c, http.StatusOK, updatedUser, nil)

}
