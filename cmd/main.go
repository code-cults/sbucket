package main

import (
	"log"
	"github.com/srirangamuc/sbucket/internal/config"
	"github.com/srirangamuc/sbucket/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/srirangamuc/sbucket/internal/api/handler"
	"github.com/srirangamuc/sbucket/internal/middleware"
)

func main(){
	config.LoadEnv()
	db.Connect()
	db.AutoMigrateModels()
	app := fiber.New()
	app.Get("/",func(c *fiber.Ctx) error{
		return c.SendString("App is live")
	})
	app.Post("/signup",handler.Signup)
	app.Post("/login",handler.Login)

	app.Post("/bucket",middleware.RequireAuth(),handler.CreateBucket)

	app.Get("/me",middleware.RequireAuth(),func(c *fiber.Ctx) error{
		userID := c.Locals("userID").(int)
		return c.JSON(fiber.Map{
			"message":"Welcome user!",
			"userID" : userID,
		})
	})

	port := config.GetEnv("PORT","3000")
	if err:=app.Listen(":"+port);err!=nil {
		log.Fatal(err)
	}
}
