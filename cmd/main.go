package main

import (
	"log"
	"github.com/srirangamuc/sbucket/internal/config"
	"github.com/srirangamuc/sbucket/internal/db"
	"github.com/srirangamuc/sbucket/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/srirangamuc/sbucket/internal/api/handler"
)

func main(){
	config.LoadEnv()
	db.Connect()
	db.DB.AutoMigrate(&model.User{})
	app := fiber.New()
	app.Get("/",func(c *fiber.Ctx) error{
		return c.SendString("App is live")
	})
	app.Post("/signup",handler.Signup)
	app.Post("/login",handler.Login)
	port := config.GetEnv("PORT","3000")
	if err:=app.Listen(":"+port);err!=nil {
		log.Fatal(err)
	}
}
