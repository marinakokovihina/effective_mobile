package http

import (
	"effective_mobile/internal/person"

	"github.com/gofiber/fiber/v2"
)

func MapRoutes(r fiber.Router, h person.HTTPhandler) {
	r.Post("/person", h.CreatePerson())
	r.Put("/person", h.UpdatePerson())
	r.Delete("/person", h.DeletePerson())
	r.Get("/persons", h.FetchPersons())
}
