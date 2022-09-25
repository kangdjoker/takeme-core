package storage

import (
	"context"
	"io"
	"mime/multipart"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/takeme-id/core/utils"
)

var StorageClient *minio.Client

func SetupStorage() error {

	endpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: true,
	})

	if err != nil {
		return err
	}

	StorageClient = minioClient

	return nil
}

func SaveFile(file io.Reader, fileHeader multipart.FileHeader) (error, string) {
	bucketName := os.Getenv("MINIO_BUCKET")
	objectName := utils.GenerateMediumCode() + fileHeader.Filename

	_, err := StorageClient.PutObject(context.Background(), bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return utils.ErrorInternalServer(utils.SaveFileFailed, err.Error()), ""
	}

	return nil, objectName
}
