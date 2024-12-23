package helper

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"simple-crud-clean-architecture/internal/entity"
	"simple-crud-clean-architecture/internal/model"
	"strings"
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

	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioRootUser, minioRootPassword, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	err = minioClient.MakeBucket(ctx, minioTicketsBucket, minio.MakeBucketOptions{Region: minioLocation})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, minioTicketsBucket)
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

func (m *Minio) getBucketName() string {
	BucketName := m.minioBucketName
	return BucketName
}

func (m *Minio) getEndpoint() string {
	Endpoint := m.enpoint
	return Endpoint
}

func (m *Minio) UploadFileToMinio(ctx context.Context, file *multipart.FileHeader, path string) (*model.MinioFileResponse, error) {
	uploadFile, err := file.Open()
	if err != nil {
		m.log.Warnf("parse file error" + err.Error())
		return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
	}
	defer uploadFile.Close()

	fileKey := path + string(RandomNumber(31)) + "_" + file.Filename
	contentType := file.Header.Get("Content-Type")

	s3PutObjectOutput, err := m.MinioClient.PutObject(ctx, m.getBucketName(), fileKey, uploadFile, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		m.log.Warnf("failed to upload file to S3" + err.Error())

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fileResponse := new(model.MinioFileResponse)
	fileResponse.ChecksumCRC32 = s3PutObjectOutput.ChecksumCRC32
	fileResponse.ChecksumCRC32C = s3PutObjectOutput.ChecksumCRC32C
	fileResponse.ChecksumSHA1 = s3PutObjectOutput.ChecksumSHA1
	fileResponse.ChecksumSHA256 = s3PutObjectOutput.ChecksumSHA256
	fileResponse.ETag = s3PutObjectOutput.ETag
	fileResponse.Expiration = s3PutObjectOutput.Expiration

	fileURL, err := m.MinioClient.PresignedGetObject(ctx, m.getBucketName(), fileKey, 1*time.Hour, nil)
	if err != nil {
		m.log.Warnf("failed to generate presigned URL:" + err.Error())

		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fileResponse.URL = fileURL.String()
	fileResponse.Filename = fileKey
	fileResponse.Mimetype = contentType
	fileResponse.Size = file.Size

	return fileResponse, nil
}

func (m *Minio) DeleteFromMinio(ctx context.Context, fileName string) (bool, error) {

	err := m.MinioClient.RemoveObject(ctx, m.getBucketName(), fileName, minio.RemoveObjectOptions{ForceDelete: true})
	if err != nil {
		return false, fmt.Errorf("failed to delete file: %w", err)
	}

	return true, nil
}

// generateAES128Key menghasilkan kunci AES-128 dan menyimpannya ke file
func generateAES128Key(filePath string) error {
	key := make([]byte, 16) // 128 bits
	if _, err := rand.Read(key); err != nil {
		return fmt.Errorf("failed to generate random key: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create key file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(key); err != nil {
		return fmt.Errorf("failed to write key to file: %w", err)
	}

	return nil
}

func (m *Minio) GenerateSignedURL(objectPath string, expiry time.Duration) (string, error) {
	ctx := context.Background()
	reqParams := make(url.Values)
	// Tambahkan parameter query jika diperlukan
	fmt.Println(objectPath)

	presignedURL, err := m.MinioClient.PresignedGetObject(ctx, m.getBucketName(), objectPath, expiry, reqParams)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (m *Minio) DownloadObject(ctx context.Context, objectPath string) ([]byte, error) {
	object, err := m.MinioClient.GetObject(ctx, m.getBucketName(), objectPath, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// UploadKey mengupload file kunci ke S3/MinIO dan mengembalikan signed URL-nya
func (m *Minio) UploadKey(ctx context.Context, videoID string, keyFilePath string) (string, error) {
	keyFilename := "encryption.key" // Nama file kunci di bucket
	s3Key := fmt.Sprintf("videos/%s/secrets/%s", videoID, keyFilename)

	// Buka file kunci
	file, err := os.Open(keyFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to open key file: %w", err)
	}
	defer file.Close()

	// Upload file kunci ke S3/MinIO
	_, err = m.MinioClient.PutObject(ctx, m.getBucketName(), s3Key, file, -1, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload key to S3: %w", err)
	}

	// Generate signed URL untuk kunci
	signedURL, err := m.GenerateSignedURL(s3Key, time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL for key: %w", err)
	}

	return signedURL, nil
}

// generateIV menghasilkan IV AES-128 dalam format hex dan menyimpannya ke file
func generateIV(filePath string) (string, error) {
	iv := make([]byte, 16) // 128 bit
	if _, err := rand.Read(iv); err != nil {
		return "", fmt.Errorf("failed to generate random IV: %w", err)
	}

	ivHex := hex.EncodeToString(iv)

	if err := os.WriteFile(filePath, []byte(ivHex), 0644); err != nil {
		return "", fmt.Errorf("failed to write IV to file: %w", err)
	}

	return ivHex, nil
}

// createKeyInfoFile membuat file key_info yang diperlukan oleh ffmpeg
func createKeyInfoFile(keyURL string, keyFilePath string, iv string, keyInfoFilePath string) error {
	content := fmt.Sprintf("%s\n%s\n%s\n", keyURL, keyFilePath, iv)

	if err := os.WriteFile(keyInfoFilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write key info file: %w", err)
	}

	return nil
}

type VideoJob struct {
	S3Client   *minio.Client
	BucketName string
	InputFile  string
	OutputFile string
}

// uploadToS3 mengupload file ke S3/MinIO dengan ContentType yang sesuai
func (m *Minio) uploadToS3(ctx context.Context, filePath string, s3Key string) error {
	// Buka file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for upload: %w", err)
	}
	defer file.Close()

	// Ambil info file
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	// Tentukan ContentType berdasarkan ekstensi file
	contentType := "application/octet-stream" // Default
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".m3u8":
		contentType = "application/vnd.apple.mpegurl"
	case ".ts":
		contentType = "video/MP2T"
	case ".key":
		contentType = "application/octet-stream"
	case ".mp4":
		contentType = "video/mp4"
	}

	// Upload file ke S3/MinIO
	_, err = m.MinioClient.PutObject(ctx, m.getBucketName(), s3Key, file, stat.Size(), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	log.Printf("Successfully uploaded %s to S3 bucket %s", s3Key, m.getBucketName())
	return nil
}

// uploadDirectoryToS3 mengupload semua file dalam direktori ke S3/MinIO dengan path relatif
func (m *Minio) uploadDirectoryToS3(ctx context.Context, dirPath string, s3Path string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(dirPath, file.Name())

			// Dapatkan path relatif dari direktori output
			relativePath, err := filepath.Rel(dirPath, filePath)
			if err != nil {
				return fmt.Errorf("failed to get relative path for %s: %w", filePath, err)
			}

			// Bentuk object key dengan menggunakan path relatif
			s3Key := fmt.Sprintf("%s/%s", s3Path, relativePath)

			// Upload file ke S3/MinIO
			err = m.uploadToS3(ctx, filePath, s3Key)
			if err != nil {
				return fmt.Errorf("failed to upload file to S3: %w", err)
			}
		}
	}

	return nil
}

func (m *Minio) ProcessVideo(ctx context.Context, inputFile string, videoID uint) (*entity.SectionVideo, error) {
	videoIDStr := fmt.Sprintf("%d", videoID) // Konversi uint ke string

	// Generate AES-128 key
	keyFilePath := fmt.Sprintf("/tmp/%s.key", videoIDStr)
	err := generateAES128Key(keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate AES-128 key: %w", err)
	}

	// Upload key ke S3/MinIO dan dapatkan signed URL
	keyURL, err := m.UploadKey(ctx, videoIDStr, keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to upload AES-128 key: %w", err)
	}

	// Generate IV yang valid
	ivFilePath := fmt.Sprintf("/tmp/%s.iv", videoIDStr)
	iv, err := generateIV(ivFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate IV: %w", err)
	}

	// Buat file key_info untuk ffmpeg dengan IV yang benar
	keyInfoFilePath := fmt.Sprintf("/tmp/%s.keyinfo", videoIDStr)
	err = createKeyInfoFile(keyURL, keyFilePath, iv, keyInfoFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create key info file: %w", err)
	}

	// Menyiapkan direktori untuk output HLS
	outputDir := "videos"
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		if err := os.Mkdir(outputDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	// Menentukan format encoding dan resolusi untuk video
	formats := map[string]struct {
		BitRate string
		Res     [2]int
	}{
		"360p":  {"500k", [2]int{480, 360}},
		"480p":  {"1000k", [2]int{858, 480}},
		"720p":  {"2000k", [2]int{1280, 720}},
		"1080p": {"4000k", [2]int{1920, 1080}},
	}

	// Menggunakan ffmpeg untuk menghasilkan HLS dengan enkripsi
	for res, format := range formats {
		resOutputDir := fmt.Sprintf("/tmp/%s", res)
		if _, err := os.Stat(resOutputDir); os.IsNotExist(err) {
			if err := os.Mkdir(resOutputDir, os.ModePerm); err != nil {
				return nil, fmt.Errorf("failed to create output directory for %s: %w", res, err)
			}
		}

		// Bentuk nama playlist HLS yang sesuai
		playlistName := fmt.Sprintf("%s.m3u8", res)

		cmd := exec.Command("ffmpeg",
			"-i", inputFile,
			"-c:v", "libx264",
			"-preset", "fast",
			"-b:v", format.BitRate,
			"-c:a", "aac",
			"-b:a", "128k",
			"-s", fmt.Sprintf("%dx%d", format.Res[0], format.Res[1]),
			"-hls_key_info_file", keyInfoFilePath, // Tambahkan opsi ini untuk enkripsi
			"-f", "hls",
			"-hls_time", "10",
			"-hls_playlist_type", "vod",
			fmt.Sprintf("%s/%s", resOutputDir, playlistName),
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("FFmpeg error for %s: %s", res, output)
			return nil, fmt.Errorf("failed to encode video to %s: %w", res, err)
		}

		log.Printf("HLS encoding for %s completed: %s/%s", res, resOutputDir, playlistName)

		// Upload direktori HLS ke S3 dengan path relatif
		s3Path := fmt.Sprintf("%s/%s/%s", outputDir, videoIDStr, res)
		err = m.uploadDirectoryToS3(ctx, resOutputDir, s3Path)
		if err != nil {
			return nil, fmt.Errorf("failed to upload HLS files to S3 for %s: %w", res, err)
		}

		// // Update progres encoding di Redis (jika diperlukan)
		// err = m.UpdateEncodingProgress(videoIDStr, res, 100.0) // Misal 100% setelah selesai
		// if err != nil {
		// 	log.Printf("Failed to update encoding progress for %s: %v", res, err)
		// }

		// Hapus direktori lokal setelah upload
		if err := os.RemoveAll(resOutputDir); err != nil {
			log.Printf("Failed to delete local directory %s: %v", resOutputDir, err)
		}
	}

	// // Update record video di database
	mediaIDName := fmt.Sprintf("%s.m3u8", videoIDStr)
	mediaDir := fmt.Sprintf("videos/%s", videoIDStr)
	// // err = UpdateVideoRecord(db, videoID, mediaDir, m.BucketName, mediaIDName)
	// // if err != nil {
	// // 	return fmt.Errorf("failed to update video record: %w", err)
	// // }

	// Hapus file kunci lokal dan file key_info setelah selesai
	if err := os.Remove(keyFilePath); err != nil {
		log.Printf("Failed to delete key file %s: %v", keyFilePath, err)
	}
	if err := os.Remove(keyInfoFilePath); err != nil {
		log.Printf("Failed to delete key info file %s: %v", keyInfoFilePath, err)
	}
	if err := os.Remove(ivFilePath); err != nil {
		log.Printf("Failed to delete IV file %s: %v", ivFilePath, err)
	}

	encodedVideoEntity := &entity.SectionVideo{
		ID:        1,
		SectionID: 1,
		MediaID:   mediaIDName,
		Bucket:    m.getBucketName(),
		Dir:       mediaDir,
	}

	return encodedVideoEntity, nil
}
