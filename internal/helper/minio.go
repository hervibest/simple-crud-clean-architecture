package helper

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Minio struct {
	MinioClient     *minio.Client
	minioBucketName string
	enpoint         string
	log             *logrus.Logger
}

func NewMinio(viper *viper.Viper, log *logrus.Logger) *Minio {
	ctx := context.Background()
	var minioClient *minio.Client
	var err error

	minioHost := viper.GetString("minio.minio_host")
	minioPort := viper.GetString("minio.minio_port")
	minioRootUser := viper.GetString("minio.minio_root_user")
	minioRootPassword := viper.GetString("minio.minio_root_password")
	minioTicketsBucket := viper.GetString("minio.minio_tickets_bucket")
	minioLocation := viper.GetString("minio.minio_location")
	endpoint := minioHost + ":" + minioPort

	minioClient, err = minio.New(endpoint, &minio.Options{ // Initialize minio client object.
		Creds:  credentials.NewStaticV4(minioRootUser, minioRootPassword, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = minioClient.MakeBucket(ctx, minioTicketsBucket, minio.MakeBucketOptions{Region: minioLocation})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, minioTicketsBucket) // Check to see if we already own this bucket (which hminioens if you run this twice)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", minioTicketsBucket)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", minioTicketsBucket)
	}

	log.Printf("Successfully connected %s\n", minioTicketsBucket)

	return &Minio{
		MinioClient:     minioClient,
		minioBucketName: minioTicketsBucket,
		log:             log,
	}
}

func (m *Minio) GetBucketName() string {
	BucketName := m.minioBucketName
	return BucketName
}

func (m *Minio) GetEnpoint() string {
	Endpoint := m.enpoint
	return Endpoint
}

func (m *Minio) UploadFileToMinio(ctx context.Context, file *multipart.FileHeader, path string) (any, error) {
	uploadFile, err := file.Open()
	if err != nil {
		m.log.Warnf("parse file error" + err.Error())
		return false, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	defer uploadFile.Close()

	fileKey := path + string(RandomNumber(31)) + "_" + file.Filename
	contentType := file.Header.Get("Content-Type")

	s3PutObjectOutput, err := m.MinioClient.PutObject(ctx, m.GetBucketName(), fileKey, uploadFile, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		m.log.Warnf("failed to upload file to S3" + err.Error())

		return false, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	resultMap := make(map[string]interface{})
	resultMap["ChecksumCRC32"] = s3PutObjectOutput.ChecksumCRC32
	resultMap["ChecksumCRC32C"] = s3PutObjectOutput.ChecksumCRC32C
	resultMap["ChecksumSHA1"] = s3PutObjectOutput.ChecksumSHA1
	resultMap["ChecksumSHA256"] = s3PutObjectOutput.ChecksumSHA256
	resultMap["ETag"] = s3PutObjectOutput.ETag
	resultMap["Expiration"] = s3PutObjectOutput.Expiration

	videoURL, err := m.MinioClient.PresignedGetObject(ctx, m.GetBucketName(), fileKey, 1*time.Hour, nil)
	if err != nil {
		m.log.Warnf("failed to generate presigned URL:" + err.Error())

		return false, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	resultMap["URL"] = videoURL.String()
	resultMap["filename"] = fileKey
	resultMap["mimetype"] = contentType
	resultMap["size"] = file.Size

	return resultMap, nil
}

func (m *Minio) DeleteFromMinio(ctx context.Context, fileName string) (bool, error) {

	err := m.MinioClient.RemoveObject(ctx, m.GetBucketName(), fileName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		return false, fmt.Errorf("failed to delete file: %w", err)
	}

	return true, nil
}
