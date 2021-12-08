package routes

import (
	"errors"
	"time"

	"github.com/anabeto93/fiber-api/database"
	"github.com/anabeto93/fiber-api/models"
	"github.com/gofiber/fiber/v2"
)

type Order struct {
	ID uint `json:"id"`
	User User `json:"user"`
	Product Product `json:"product"`
	CreatedAt time.Time `json:"order_date"`
}

func CreateResponseOrder(order models.Order, user User, product Product) Order {
	return Order{ID: order.ID, User: user, Product: product, CreatedAt: order.CreatedAt}
}

func CreateOrder(ctx *fiber.Ctx) error {
	var order models.Order

	if err := ctx.BodyParser(&order); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	var user models.User
	if err := findUser(order.UserRefer, &user); err != nil {
		return ctx.Status(404).JSON(err.Error())
	}

	var product models.Product
	if err := findProduct(order.ProductRefer, &product); err != nil {
		return ctx.Status(404).JSON(err.Error())
	}

	database.Database.Db.Create(&order)

	resUser := CreateResponseUser(user)
	resProduct := CreateResponseProduct(product)
	response := CreateResponseOrder(order, resUser, resProduct)

	return ctx.Status(201).JSON(response)
}

func GetOrders(ctx *fiber.Ctx) error {
	orders := []models.Order{}
	database.Database.Db.Find(&orders)

	response := []Order{}

	for _, order := range orders {
		// this is inefficient but don't think too much about this
		var user models.User
		var product models.Product
		findUser(order.UserRefer, &user)
		findProduct(order.ProductRefer, &product)
		
		resUser := CreateResponseUser(user)
		resProduct := CreateResponseProduct(product)

		tempOrder := CreateResponseOrder(order, resUser, resProduct)
		response = append(response, tempOrder)
	}

	return ctx.Status(200).JSON(response)
}

func findOrder(id int, order *models.Order) error {
	database.Database.Db.Find(&order, "id = ?", id)

	if order.ID == 0 {
		return errors.New("Order does not exist")
	}

	return nil
}

func GetOrder(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id"); if err != nil {
		return ctx.Status(400).JSON(":id must be an integer")
	}

	var order models.Order

	if err := findOrder(id, &order); err != nil {
		return ctx.Status(404).JSON(err.Error())
	}

	var user models.User
	var product models.Product
	findUser(order.UserRefer, &user)
	findProduct(order.ProductRefer, &product)

	userRes := CreateResponseUser(user)
	productRes := CreateResponseProduct(product)

	response := CreateResponseOrder(order, userRes, productRes)
	return ctx.Status(200).JSON(response)
}

func UpdateOrder(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id"); if err != nil {
		return ctx.Status(400).JSON(":id must be an integer")
	}

	type UpdateOrder struct {
		UserId int `json:"user_id"`
		ProductId int `json:"product_id"`
	}

	var update UpdateOrder

	if err := ctx.BodyParser(&update); err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	var order models.Order

	if err := findOrder(id, &order); err != nil {
		return ctx.Status(404).JSON(err.Error())
	}

	order.UserRefer = update.UserId
	order.ProductRefer = update.ProductId

	database.Database.Db.Save(&order)

	var user models.User
	var product models.Product
	findUser(order.UserRefer, &user)
	findProduct(order.ProductRefer, &product)

	userRes := CreateResponseUser(user)
	productRes := CreateResponseProduct(product)

	response := CreateResponseOrder(order, userRes, productRes)

	return ctx.Status(200).JSON(response)
}

func DeleteOrder(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id"); if err != nil {
		return ctx.Status(400).JSON(":id must be an integer")
	}

	var order models.Order

	if err := findOrder(id, &order); err != nil {
		return ctx.Status(404).JSON(err.Error())
	}

	if err := database.Database.Db.Delete(&order).Error; err != nil {
		return ctx.Status(400).JSON(err.Error())
	}

	return ctx.Status(200).SendString("Order deleted.")
}