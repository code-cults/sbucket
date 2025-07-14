package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/srirangamuc/sbucket/internal/db"
	"github.com/srirangamuc/sbucket/internal/model"
	"github.com/srirangamuc/sbucket/internal/storage"
	"github.com/google/uuid"
)

func GenerateSignedDownloadURL(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	bucketIDstr := c.Params("bucketID")
	fileName := c.Params("fileName")

	bucketID, err := uuid.Parse(bucketIDstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error" : "Invalid Bucket ID",
		})
	}

	var bucket model.Bucket
	if err := db.DB.First(&bucket,"id = ? AND owner_id = ?",bucketID,userID).Error; err != nil {
		return c.Status(403).JSON(fiber.Map{
			"error" : "Unauthorized or Bucket not Found",
		})
	}

	var fileMeta model.File
	if err := db.DB.First(&fileMeta,"bucket_id = ? AND file_name = ?",bucketID,fileName).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error" : "File not found",
		})
	}

	objectKey := fmt.Sprintf("%s/%s",bucketID.String(),fileName)
	url, err := storage.GeneratePresignedURL("sbucket", objectKey, time.Minute*10)
	if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to generate signed URL"})
    }

    return c.JSON(fiber.Map{
        "url": url,
        "expires_in": "10m",
    })
}

func GenerateSignedUploadURL(c *fiber.Ctx) error {
    userID := c.Locals("userID").(int)

    bucketIDStr := c.Params("bucketID")
    fileName := c.Params("filename")

    bucketID, err := uuid.Parse(bucketIDStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid bucket ID"})
    }
    var bucket model.Bucket
    if err := db.DB.First(&bucket, "id = ? AND owner_id = ?", bucketID, userID).Error; err != nil {
        return c.Status(403).JSON(fiber.Map{"error": "unauthorized or bucket not found"})
    }

    objectKey := fmt.Sprintf("%s/%s", bucketID.String(), fileName)

    url, err := storage.GeneratePresignedPutURL("nebulabucket", objectKey, time.Minute*10)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to generate upload URL"})
    }

    return c.JSON(fiber.Map{
        "upload_url": url,
        "expires_in": "10m",
    })
}