package storage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"log"
)

var Client *minio.Client

func InitMinIO() {
	var err error
	Client, err = minio.New("localhost:9000", &minio.Options{
		Creds: credentials.NewStaticV4("minio_user", "minio_password", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal("Failed to Initialize MinIO")
	}
	log.Println("MinIO is Initialized successfully")
}

func UploadToMinIO(bucketName, objectName string, file io.Reader, fileSize int64, contentType string) error {
	_, err := Client.PutObject(context.Background(), bucketName, objectName, file, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func DownloadFromMinIO(bucketName,objectName string) (io.Reader,error) {
	obj, err := Client.GetObject(context.Background(),bucketName,objectName,minio.GetObjectOptions{})
	if err!= nil {
		return nil, err
	}
	return obj, nil
}