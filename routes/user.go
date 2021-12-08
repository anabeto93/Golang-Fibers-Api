package routes

import (
	"errors"

	"github.com/anabeto93/fiber-api/database"
	"github.com/anabeto93/fiber-api/models"
	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID uint `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}

func CreateResponseUser(user models.User) User {
	return User{ID: user.ID, FirstName: user.FirstName, LastName: user.LastName}
}

func CreateUser(ctx *fiber.Ctx) error {
	var user models.User

	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	database.Database.Db.Create(&user)
	response := CreateResponseUser(user)

	return ctx.Status(201).JSON(response)
}

func GetUsers(ctx *fiber.Ctx) error {
	users := []models.User{}

	database.Database.Db.Find(&users)
	responseUsers := []User{}

	for _, user := range users {
		tempUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, tempUser)
	}

	return ctx.Status(200).JSON(responseUsers)
}

func findUser(id int, user *models.User) error {
	database.Database.Db.Find(&user, "id = ?", id)
	if user.ID == 0 {
		return errors.New("User does not exist")
	}
	return nil
}

func GetUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id"); if err != nil {
		return ctx.Status(400).JSON(":id must be an integer")
	}

	var user models.User

	if err := findUser(id, &user); err != nil {
		return ctx.Status(404).JSON(err.Error())
	}

	response := CreateResponseUser(user)
	
	return ctx.Status(200).JSON(response)
}

func UpdateUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id"); if err != nil {
		return ctx.Status(400).JSON(":id must be an integer")
	}

	type UpdateUser struct {
		FirstName string `json:"first_name"`
		LastName string `json:"last_name"`
	}

	var update UpdateUser

	if err := ctx.BodyParser(&update); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	var user models.User

	if err := findUser(id, &user); err != nil {
		return ctx.Status(404).JSON(err.Error())
	}

	user.FirstName = update.FirstName
	user.LastName = update.LastName

	database.Database.Db.Save(&user)

	response := CreateResponseUser(user)
	return ctx.Status(200).JSON(response)
}

func DeleteUser(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id"); if err != nil {
		return ctx.Status(400).JSON(":id must be an integer")
	}

	var user models.User

	if err := findUser(id, &user); err != nil {
		return ctx.Status(404).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&user).Error; err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	return ctx.Status(200).SendString("User deleted.")
}