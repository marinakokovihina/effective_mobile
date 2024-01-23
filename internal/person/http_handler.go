package person

import "github.com/gofiber/fiber/v2"

type HTTPhandler interface {
	CreatePerson() fiber.Handler
	UpdatePerson() fiber.Handler
	DeletePerson() fiber.Handler
	FetchPersons() fiber.Handler
}
