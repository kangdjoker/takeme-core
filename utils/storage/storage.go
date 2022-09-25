package storage

import (
	"context"
	"io"
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

func SaveFile(bucketName string, objectName string, file io.Reader, fileSize int64) error {
	_, err := StorageClient.PutObject(context.Background(), bucketName, objectName, file, fileSize, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return utils.ErrorInternalServer(utils.SaveFileFailed, err.Error())
	}

	return nil
}
