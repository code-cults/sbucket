package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/srirangamuc/sbucket/internal/service"
)

type CreateBucketRequest struct {
	Name string `json:"name"`
}

func CreateBucket(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)
	var req CreateBucketRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Invalid Request",
		})
	}
	bucket, err := service.CreateBucket(userID,req.Name)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message" : "Bucket Created Successfully",
		"bucket_id" : bucket.ID,
	})
}