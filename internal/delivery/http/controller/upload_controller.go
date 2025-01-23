package http

import (
	"net/url"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UploadController struct {
	Log   *logrus.Logger
	Minio *helper.Minio
}

func NewUploadController(log *logrus.Logger, minio helper.Minio) *UploadController {
	return &UploadController{
		Log:   log,
		Minio: &minio,
	}
}

func (u *UploadController) UploadFile(ctx *fiber.Ctx) error {

	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}

	const maxFileSize = 0.5 * 1024 * 1024 // 5MB

	if file.Size > maxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "File size exceeds the 5MB limit")
	}

	upload, err := u.Minio.UploadFileToMinio(ctx.UserContext(), file, "avatar")
	if err != nil {
		return err
	}

	return ctx.JSON(model.DataResponse[any]{
		Success: true,
		Data:    upload,
	})
}

func (c *UploadController) DeleteFile(ctx *fiber.Ctx) error {
	fileName := ctx.Params("fileName")

	decodedFileName, err := url.QueryUnescape(fileName)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "failed to decode filename")
	}

	valid, err := c.Minio.DeleteFromMinio(ctx.UserContext(), decodedFileName)
	if err != nil {
		c.Log.WithError(err).Error("error deleting uploaded file in minio")
		return fiber.ErrNotFound
	}

	return ctx.JSON(model.DataResponse[any]{
		Success: true,
		Data:    valid,
	})
}
