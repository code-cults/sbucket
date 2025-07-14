package main

import (
	"log"
	"github.com/srirangamuc/sbucket/internal/config"
	"github.com/srirangamuc/sbucket/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/srirangamuc/sbucket/internal/api/handler"
	"github.com/srirangamuc/sbucket/internal/middleware"
	"github.com/srirangamuc/sbucket/internal/storage"
)

func main(){
	config.LoadEnv()
	db.Connect()
	db.AutoMigrateModels()
	storage.InitMinIO()
	app := fiber.New()

	api := app.Group("/api")
	protected := api.Use(middleware.RequireAuth())


	app.Get("/",func(c *fiber.Ctx) error{
		return c.SendString("App is live")
	})
	app.Post("/signup",handler.Signup)
	app.Post("/login",handler.Login)

	protected.Post("/bucket",handler.CreateBucket)
	protected.Post("/bucket/:bucketID/upload",handler.UploadFile)
	protected.Get("/bucket/:bucketID/file/:fileName",handler.DownloadFile)
	protected.Get("/bucket/:bucketID/files",handler.ListOfFilesInBucket)
	protected.Get("/bucket/:bucketID/file/:fileName/signed",handler.GenerateSignedDownloadURL)
	protected.Get("/bucket/:bucketID/upload-url/:fileName",handler.GeneratedSignedUploadURL)

	port := config.GetEnv("PORT","3000")
	if err:=app.Listen(":"+port);err!=nil {
		log.Fatal(err)
	}
}
