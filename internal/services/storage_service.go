// internal/services/storage_service.go
package services

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageService struct {
    minioClient *minio.Client
    bucketName  string
}

func NewStorageService() (*StorageService, error) {
    log.Printf("Initializing MinIO service...")
    
    // Configuration MinIO
    endpoint := "minio.okloud-hub.com"
    accessKey := "minio"
    secretKey := "minio123"
    bucketName := "mymusics"

    log.Printf("MinIO Config - Endpoint: %s, Bucket: %s", endpoint, bucketName)

    // Créer le client
    minioClient, err := minio.New(endpoint, &minio.Options{
        Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
        Secure: true,
        Region: "", // Enlever la région pour tester
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{
                InsecureSkipVerify: true, // Pour le développement
            },
        },
    })
    
    if err != nil {
        log.Printf("Failed to create MinIO client: %v", err)
        return nil, fmt.Errorf("erreur création client minio: %v", err)
    }

    log.Printf("MinIO client created successfully")

    // Vérifier la connexion sans tester le bucket tout de suite
    return &StorageService{
        minioClient: minioClient,
        bucketName:  bucketName,
    }, nil
}

// Ajouter une méthode de test séparée
func (s *StorageService) TestConnection() error {
    exists, err := s.minioClient.BucketExists(context.Background(), s.bucketName)
    if err != nil {
        return fmt.Errorf("erreur vérification bucket: %v", err)
    }
    log.Printf("Bucket %s exists: %v", s.bucketName, exists)
    return nil
}

func (s *StorageService) ListFiles() ([]minio.ObjectInfo, error) {
    ctx := context.Background()
    
    // Utilisez AWS V4 signature
    opts := minio.ListObjectsOptions{
        Recursive: true,
        UseV1: false,  // Force l'utilisation de V2
    }
    
    objectCh := s.minioClient.ListObjects(ctx, s.bucketName, opts)
    var objects []minio.ObjectInfo

    for object := range objectCh {
        if object.Err != nil {
            return nil, fmt.Errorf("erreur listage fichiers: %v", object.Err)
        }
        objects = append(objects, object)
    }

    return objects, nil
}

func (s *StorageService) UploadFile(file *multipart.FileHeader) (string, error) {
    src, err := file.Open()
    if err != nil {
        return "", fmt.Errorf("erreur ouverture fichier: %v", err)
    }
    defer src.Close()

    objectName := uuid.New().String() + filepath.Ext(file.Filename)
    
    // Upload avec metadata
    _, err = s.minioClient.PutObject(
        context.Background(),
        s.bucketName,
        objectName,
        src,
        file.Size,
        minio.PutObjectOptions{
            ContentType: file.Header.Get("Content-Type"),
            UserMetadata: map[string]string{
                "originalname": file.Filename,
            },
        },
    )
    if err != nil {
        return "", fmt.Errorf("erreur upload: %v", err)
    }

    return objectName, nil
}

// Ajoutez une méthode pour récupérer un fichier
func (s *StorageService) GetFile(objectName string) (io.Reader, error) {
    obj, err := s.minioClient.GetObject(
        context.Background(),
        s.bucketName,
        objectName,
        minio.GetObjectOptions{},
    )
    if err != nil {
        return nil, fmt.Errorf("erreur récupération fichier: %v", err)
    }

    return obj, nil
}

// Méthode pour obtenir l'URL signée temporaire (utile pour le streaming)
func (s *StorageService) GetPresignedURL(objectName string, expiry time.Duration) (string, error) {
    presignedURL, err := s.minioClient.PresignedGetObject(
        context.Background(),
        s.bucketName,
        objectName,
        expiry,
        nil,
    )
    if err != nil {
        return "", fmt.Errorf("erreur génération URL: %v", err)
    }

    return presignedURL.String(), nil
}