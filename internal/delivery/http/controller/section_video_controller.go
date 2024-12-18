package http

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"simple-crud-clean-architecture/internal/helper"
	"simple-crud-clean-architecture/internal/model"
	"simple-crud-clean-architecture/internal/usecase"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SectionVideoController struct {
	UseCase   *usecase.SecVideoUseCase
	Log       *logrus.Logger
	Validator helper.CustomValidator
	Minio     *helper.Minio
}

func NewSecVideoController(useCase *usecase.SecVideoUseCase, log *logrus.Logger, validator helper.CustomValidator, minio *helper.Minio) *SectionVideoController {
	return &SectionVideoController{
		UseCase:   useCase,
		Log:       log,
		Validator: validator,
		Minio:     minio,
	}
}

func (c *SectionVideoController) Create(ctx *fiber.Ctx) error {

	request := new(model.CreateSecVideoRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.Create(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error creating section video")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.SecVideoResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *SectionVideoController) List(ctx *fiber.Ctx) error {

	request := &model.SearchSecVideosRequest{
		Title: ctx.Query("title", ""),
		Notes: ctx.Query("notes", ""),
		Page:  ctx.QueryInt("page", 1),
		Size:  ctx.QueryInt("size", 10),
	}

	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	responses, pageMetadata, err := c.UseCase.Search(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error searching section videos")
		return err
	}

	baseURL := ctx.BaseURL() + ctx.Path()
	helper.GeneratePageURLs(baseURL, pageMetadata)

	return ctx.JSON(model.DataResponse[[]model.SecVideoResponse]{Success: true,
		Data:   responses,
		Paging: pageMetadata,
	})

}

func (c *SectionVideoController) Get(ctx *fiber.Ctx) error {
	parsedUUID, err := uuid.Parse(ctx.Params("secVideoID"))

	c.Log.Infof(parsedUUID.String())

	if err != nil {
		c.Log.WithError(err).Error("error parsing section video controller")
		return fiber.ErrBadRequest
	}

	request := &model.GetSecVideoRequest{
		UUID: parsedUUID,
	}

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error getting section video")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.SecVideoResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *SectionVideoController) Update(ctx *fiber.Ctx) error {

	request := new(model.UpdateSecVideoRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("secVideoID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing section video uuid")
		return fiber.ErrBadRequest
	}

	request.VideoUUID = parsedUUID

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.Update(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error updating section video")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.SecVideoResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *SectionVideoController) Delete(ctx *fiber.Ctx) error {

	request := new(model.DeleteSecVideoRequest)
	if err := ctx.BodyParser(request); err != nil {
		c.Log.WithError(err).Error("error parsing request body")
		return fiber.ErrBadRequest
	}

	parsedUUID, err := uuid.Parse(ctx.Params("secVideoID"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing section video uuid")
		return fiber.ErrBadRequest
	}

	request.VideoUUID = parsedUUID

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.Delete(ctx.UserContext(), request)
	if err != nil {
		c.Log.WithError(err).Error("error deleting secttion video")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.SecVideoResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *SectionVideoController) UploadVideo(ctx *fiber.Ctx) error {

	parsedSecVideoUUID, err := uuid.Parse(ctx.Params("secVideoId"))
	if err != nil {
		c.Log.WithError(err).Error("error parsing section video uuid")
		return fiber.ErrBadRequest
	}

	sectionUUID := ctx.FormValue("SectionUUID")
	if sectionUUID == "" {
		c.Log.Error("missing SectionUUID in form data")
		return fiber.NewError(fiber.StatusBadRequest, "missing SectionUUID")
	}

	parsedSectionUUID, err := uuid.Parse(sectionUUID)
	if err != nil {
		c.Log.WithError(err).Error("error parsing SectionUUID")
		return fiber.ErrBadRequest
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusUnprocessableEntity, "missing file: "+err.Error())
	}

	const maxFileSize = 10 * 1024 * 1024 // 10MB
	if file.Size > maxFileSize {
		return fiber.NewError(fiber.StatusRequestEntityTooLarge, "File size exceeds the 10MB limit")
	}

	request := &model.UploadVideoRequest{
		VideoUUID:   parsedSecVideoUUID,
		SectionUUID: parsedSectionUUID,
	}

	if validationErr := c.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := c.UseCase.UploadVideo(ctx.UserContext(), file, request)
	if err != nil {
		c.Log.WithError(err).Error("error updating section video")
		return err
	}

	return ctx.JSON(model.DataResponse[*model.SecVideoResponse]{
		Success: true,
		Data:    response,
	})
}

func (c *SectionVideoController) EncodeVideo(ctx *fiber.Ctx) error {

	go func() {
		err := c.UseCase.EncodeVideos(ctx.UserContext())
		if err != nil {
			c.Log.Warnf("Failed to encode video : %+v", err)
		}

	}()

	return ctx.JSON(model.WebResponse{
		Success: true,
		Data:    nil})

}

func (s *SectionVideoController) ServeHLSPlaylist(ctx *fiber.Ctx) error {
	// videoID := ctx.Params("videoID")
	playlist := ctx.Params("playlist") // e.g., "master.m3u8"

	parsedUUID, err := uuid.Parse(ctx.Params("videoID"))

	s.Log.Infof(parsedUUID.String())

	if err != nil {
		s.Log.WithError(err).Error("error parsing section video controller")
		return fiber.ErrBadRequest
	}

	request := &model.GetSecVideoRequest{
		UUID: parsedUUID,
	}

	if validationErr := s.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := s.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		s.Log.WithError(err).Error("error getting section video")
		return err
	}

	fmt.Println(response.MediaID)
	objectPath := fmt.Sprintf("%s/1080p/%s", response.Dir, playlist)

	fmt.Println(objectPath)

	playlistData, err := s.Minio.DownloadObject(context.Background(), objectPath)
	if err != nil {
		log.Printf("Failed to download playlist: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to download playlist")
	}

	playlistContent := string(playlistData)

	modifiedPlaylist, err := s.modifyPlaylist(response, playlistContent)
	if err != nil {
		log.Printf("Failed to modify playlist: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to modify playlist")
	}

	// Set Content-Type sesuai HLS
	ctx.Set("Content-Type", "application/vnd.apple.mpegurl")
	return ctx.SendString(modifiedPlaylist)
}

func (s *SectionVideoController) ServeHLSKey(ctx *fiber.Ctx) error {
	key := ctx.Params("key") // e.g., "encryption.key"

	parsedUUID, err := uuid.Parse(ctx.Params("videoID"))

	s.Log.Infof(parsedUUID.String())

	if err != nil {
		s.Log.WithError(err).Error("error parsing section video controller")
		return fiber.ErrBadRequest
	}

	request := &model.GetSecVideoRequest{
		UUID: parsedUUID,
	}

	if validationErr := s.Validator.Validate(request); validationErr != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(model.ValidationErrorResponse{
			Success: false,
			Errors:  validationErr,
			Message: "validation error",
		})
	}

	response, err := s.UseCase.Get(ctx.UserContext(), request)
	if err != nil {
		s.Log.WithError(err).Error("error getting section video")
		return err
	}

	keyPath := fmt.Sprintf("%s/secrets/%s", response.Dir, key)

	keyData, err := s.Minio.DownloadObject(context.Background(), keyPath)
	if err != nil {
		log.Printf("Failed to download key: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to download key")
	}

	// Set Content-Type sesuai kunci enkripsi
	ctx.Set("Content-Type", "application/octet-stream")
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", key))

	return ctx.Send(keyData)
}

// modifyPlaylist mengganti URL segmen dan kunci dalam playlist dengan signed URLs
func (s *SectionVideoController) modifyPlaylist(video *model.SecVideoResponse, playlist string) (string, error) {
	lines := strings.Split(playlist, "\n")
	var modifiedLines []string

	for _, line := range lines {
		// Ganti URL kunci
		if strings.HasPrefix(line, "#EXT-X-KEY") {
			// Ekstrak key URI
			keyURI := extractKeyURI(line)
			if keyURI == "" {
				// Jika tidak ada key URI, tambahkan key info
				// Sesuaikan dengan kebutuhan Anda
				modifiedLines = append(modifiedLines, line)
				continue
			}

			// Generate URL untuk kunci
			keyURL := fmt.Sprintf("http://localhost:5000/api/video/%s/key/%s", video.UUID, filepath.Base(keyURI))
			line = fmt.Sprintf("#EXT-X-KEY:METHOD=AES-128,URI=\"%s\"", keyURL)
			modifiedLines = append(modifiedLines, line)
			continue
		}

		// Ganti URL segmen (.ts)
		if strings.HasSuffix(line, ".ts") {
			// Generate signed URL untuk segmen
			segPath := fmt.Sprintf("%s/1080p/%s", video.Dir, line)
			signedSegURL, err := s.Minio.GenerateSignedURL(segPath, time.Hour)
			if err != nil {
				log.Printf("Failed to generate signed URL for segment %s: %v", segPath, err)
				// Ganti dengan URL asli jika gagal
				modifiedLines = append(modifiedLines, line)
				continue
			}
			line = signedSegURL
			modifiedLines = append(modifiedLines, line)
			continue
		}

		// Ganti URL playlist (master playlist)
		if strings.HasSuffix(line, ".m3u8") && !strings.HasPrefix(line, "#") {
			// Generate signed URL untuk playlist
			playlistPath := fmt.Sprintf("%s/%s", video.Dir, line)
			signedPlaylistURL, err := s.Minio.GenerateSignedURL(playlistPath, time.Hour)
			if err != nil {
				log.Printf("Failed to generate signed URL for playlist %s: %v", playlistPath, err)
				// Ganti dengan URL asli jika gagal
				modifiedLines = append(modifiedLines, line)
				continue
			}
			line = signedPlaylistURL
			modifiedLines = append(modifiedLines, line)
			continue
		}

		// Line lainnya tetap tidak diubah
		modifiedLines = append(modifiedLines, line)
	}

	return strings.Join(modifiedLines, "\n"), nil
}

// extractKeyURI mengekstrak URI kunci dari baris #EXT-X-KEY
func extractKeyURI(line string) string {
	// Contoh line: #EXT-X-KEY:METHOD=AES-128,URI="http://example.com/key.key"
	parts := strings.Split(line, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "URI=") {
			uri := strings.TrimPrefix(part, "URI=")
			uri = strings.Trim(uri, "\"")
			return uri
		}
	}
	return ""
}
