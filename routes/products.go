package routes

import (
	"errors"

	"github.com/anabeto93/fiber-api/database"
	"github.com/anabeto93/fiber-api/models"
	"github.com/gofiber/fiber/v2"
)

type Product struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	SerialNumber string `json:"serial_number"`
}

func CreateResponseProduct(productModel models.Product) Product {
	return Product{ID: productModel.ID, Name: productModel.Name, SerialNumber: productModel.SerialNumber}
}

func CreateProduct(ctx *fiber.Ctx) error {
	var product models.Product

	if err := ctx.BodyParser(&product); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	database.Database.Db.Create(&product)
	
	response := CreateResponseProduct(product)

	return ctx.Status(201).JSON(response)
}

func GetProducts(ctx *fiber.Ctx) error {
	products := []models.Product{}
	database.Database.Db.Find(&products)

	response := []Product{}

	for _, product := range products {
		tempProduct := CreateResponseProduct(product)
		response = append(response, tempProduct)
	}

	return ctx.Status(200).JSON(response)
}

func findProduct(id int, product *models.Product) error {
	database.Database.Db.Find(&product, "id = ?", id)
	if product.ID == 0 {
		return errors.New("Product does not exist")
	}
	return nil
}

func GetProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id"); if err != nil {
		return ctx.Status(400).JSON(":id must be an integer")
	}

	var product models.Product

	if err := findProduct(id, &product); err != nil {
		return ctx.Status(404).JSON(err.Error())
	}

	response := CreateResponseProduct(product)
	return ctx.Status(200).JSON(response)
}

func UpdateProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id"); if err != nil {
		return ctx.Status(400).JSON(":id must be an integer")
	}

	type UpdateProduct struct {
		Name string `json:"name"`
		SerialNumber string `json:"serial_number"`
	}

	var update UpdateProduct

	if err := ctx.BodyParser(&update); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	var product models.Product

	if err := findProduct(id, &product); err != nil {
		return ctx.Status(404).JSON(err.Error())
	}

	product.Name = update.Name
	product.SerialNumber = update.SerialNumber

	database.Database.Db.Save(&product)

	response := CreateResponseProduct(product)

	return ctx.Status(200).JSON(response)
}



func DeleteProduct(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id"); if err != nil {
		return ctx.Status(400).JSON(":id must be an integer")
	}

	var product models.Product

	if err := findProduct(id, &product); err != nil {
		return ctx.Status(404).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&product).Error; err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	return ctx.Status(200).SendString("Product deleted.")
}