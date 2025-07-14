package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/srirangamuc/sbucket/internal/db"
	"github.com/srirangamuc/sbucket/internal/model"
)

func ListOfFilesInBucket(c *fiber.Ctx) error {
    userID := c.Locals("userID").(int)

    bucketIDStr := c.Params("bucketID")
    bucketID, err := uuid.Parse(bucketIDStr)
    if err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "invalid bucket ID"})
    }
    var bucket model.Bucket
    if err := db.DB.First(&bucket, "id = ? AND owner_id = ?", bucketID, userID).Error; err != nil {
        return c.Status(404).JSON(fiber.Map{"error": "bucket not found or unauthorized"})
    }
    var files []model.File
    if err := db.DB.Where("bucket_id = ?", bucketID).Find(&files).Error; err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "failed to fetch files"})
    }

    return c.JSON(fiber.Map{
        "bucket_id": bucketID,
        "files":     files,
    })
}