// internal/services/storage_service.go
package services

import (
	"context"
	"mime/multipart"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageService struct {
    minioClient *minio.Client
    bucketName  string
}

func NewStorageService() (*StorageService, error) {
    client, err := minio.New("minio:9000", &minio.Options{
        Creds:  credentials.NewStaticV4("minio", "minio123", ""),
        Secure: false,
    })
    if err != nil {
        return nil, err
    }

    bucketName := "music"
    exists, err := client.BucketExists(context.Background(), bucketName)
    if err != nil {
        return nil, err
    }

    if !exists {
        err = client.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
        if err != nil {
            return nil, err
        }
    }

    return &StorageService{
        minioClient: client,
        bucketName:  bucketName,
    }, nil
}

func (s *StorageService) UploadFile(file *multipart.FileHeader) (string, error) {
    src, err := file.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()

    objectName := uuid.New().String() + filepath.Ext(file.Filename)
    _, err = s.minioClient.PutObject(context.Background(), s.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
        ContentType: file.Header.Get("Content-Type"),
    })
    if err != nil {
        return "", err
    }

    return objectName, nil
}