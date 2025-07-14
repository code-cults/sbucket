package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/srirangamuc/sbucket/internal/db"
	"github.com/srirangamuc/sbucket/internal/model"
	"github.com/srirangamuc/sbucket/internal/storage"
)

func DownloadFile(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	bucketIDstr := c.Params("bucketID")
	fileName := c.Params("filename")

	bucketID, err := uuid.Parse(bucketIDstr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error" : "Invalid Bucket ID",
		})
	}

	var bucket model.Bucket
	if err := db.DB.First(&bucket,"id = ? AND owner_id = ?",bucketID,userID).Error; err != nil{
		return c.Status(404).JSON(fiber.Map{
			"error" : "Bucked not found or unauthorized",
		})
	}

	var fileMeta model.File
	if err := db.DB.First(&fileMeta,"bucket_id = ? AND file_name = ?",bucketID,fileName).Error; err != nil{
		return c.Status(404).JSON(fiber.Map{
			"error" : "File not found",
		})
	}

	objectKey := fmt.Sprintf("%s/%s",bucketID.String(),fileName)
	reader, err := storage.DownloadFromMinIO("sbucket",objectKey)
	if err != nil{
		return c.Status(500).JSON(fiber.Map{
			"error" : "Failed to download file from storage",
		})
	}

	c.Set("Content-Type",fileMeta.MimeType)
	c.Set("Content-Disposition" , fmt.Sprintf("attachment; filename=\"%s\"",fileName))
	return c.Status(200).SendStream(reader)


}