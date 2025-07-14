package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"time"
	"bytes"
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/srirangamuc/sbucket/internal/db"
	"github.com/srirangamuc/sbucket/internal/model"
	"github.com/srirangamuc/sbucket/internal/storage"
	"github.com/minio/minio-go/v7"
)

func UploadFile(c *fiber.Ctx) error {
	bucketIDStr := c.Params("bucketID")
	bucketID, err := uuid.Parse(bucketIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid Bucket ID",
		})
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "File is required",
		})
	}

	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to open the file",
		})
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to read file",
		})
	}

	// Calculate hash in goroutine
	hashChan := make(chan string)
	go func(data []byte) {
		hasher := sha256.New()
		hasher.Write(data)
		hashChan <- hex.EncodeToString(hasher.Sum(nil))
	}(fileBytes)

	objectName := fmt.Sprintf("%s/%s", bucketID.String(), fileHeader.Filename)

	err = storage.Client.MakeBucket(context.Background(), "sbucket", minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := storage.Client.BucketExists(context.Background(), "sbucket")
		if errBucketExists != nil || !exists {
			return c.Status(500).JSON(fiber.Map{"error": "Bucket creation failed"})
		}
	}


	// Replace 'storage' with your actual storage package variable
	err = storage.UploadToMinIO("sbucket", objectName,
		io.NopCloser(bytes.NewReader(fileBytes)),
		int64(len(fileBytes)),
		fileHeader.Header.Get("Content-Type"),
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Upload to MinIO failed",
		})
	}

	hash := <-hashChan

	fileMeta := model.File{
		BucketID:  bucketID,
		FileName:  fileHeader.Filename,
		Size:      int64(len(fileBytes)),
		Hash:      hash,
		MimeType:  fileHeader.Header.Get("Content-Type"),
		CreatedAt: time.Now(),
	}

	if err := db.DB.Create(&fileMeta).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "db insert failed"})
	}

	return c.JSON(fiber.Map{
		"message":  "upload successful",
		"file_id":  fileMeta.ID,
		"filename": fileMeta.FileName,
		"hash":     hash,
	})
}