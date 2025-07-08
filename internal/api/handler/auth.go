package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/srirangamuc/sbucket/internal/service"
)

type SignUpRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

func Signup(c *fiber.Ctx) error{
	var req SignUpRequest
	if err:= c.BodyParser(&req); err!=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Request",
		})
	}
	err := service.SignUpUser(req.Email, req.Password)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message" : "user created successfully",
	})
}

func Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err:= c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Invalid Request",
		})
	}
	token,err := service.LoginUser(req.Email,req.Password)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error" : "Invalid Email or Password",
		})
	}

	return c.JSON(fiber.Map{
		"token":token,
	})

}