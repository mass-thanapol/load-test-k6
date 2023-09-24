package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Token struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

func main() {

	app := fiber.New()

	dsn := "root:password@tcp(localhost:3306)/load_test?parseTime=True"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	db.AutoMigrate(&User{}, &Token{})

	app.Use(func(c *fiber.Ctx) error {
		fmt.Println("Request: ", c.Method(), c.Path(), string(c.Body()))
		if err := c.Next(); err != nil {
			return err
		}
		fmt.Println("Response: ", c.Method(), c.Path(), string(c.Response().Body()))
		return nil
	})

	app.Get("/generateToken", func(c *fiber.Ctx) error {
		tokenValue := uuid.New().String()
		token := Token{
			Token: tokenValue,
		}
		db.Create(&token)
		return c.JSON(fiber.Map{
			"token": tokenValue,
		})
	})

	app.Get("/getFirstUser", func(c *fiber.Ctx) error {
		var users []User
		db.First(&users)
		return c.JSON(users[0])
	})

	app.Post("/createUser", func(c *fiber.Ctx) error {
		var newUser User
		if err := c.BodyParser(&newUser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Invalid request body",
			})
		}
		db.Create(&newUser)
		return c.JSON(newUser)
	})

	app.Post("/deleteUser", func(c *fiber.Ctx) error {
		var requestBody struct {
			UserID int    `json:"userId"`
			Token  string `json:"token"`
		}
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "FAILURE",
				"message": "Invalid request body",
			})
		}
		var token Token
		if err := db.Where("token = ?", requestBody.Token).First(&token).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "FAILURE",
				"message": "Invalid token",
			})
		}
		result := db.Delete(&User{}, requestBody.UserID)
		if result.Error != nil || result.RowsAffected == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "FAILURE",
				"message": "User not found or deletion failed",
			})
		}
		return c.JSON(fiber.Map{
			"status": "SUCCESS",
		})
	})

	app.Listen(":3000")
}
