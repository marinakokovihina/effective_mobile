package http

import (
	"effective_mobile/internal/person"
	"effective_mobile/internal/person/model"
	"effective_mobile/status"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type handler struct {
	uc person.UC
}

func NewHandler(uc person.UC) person.HTTPhandler {
	return &handler{
		uc: uc,
	}
}

func (h handler) CreatePerson() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.CreatePersonParams
		if err := ctx.BodyParser(&params); err != nil {
			return status.ErrInput.Wrap(fmt.Errorf("cannot parse input body"))
		}
		err := h.uc.CreatePerson(ctx.Context(), params)
		if err != nil {
			return err
		}
		return ctx.JSON(status.HTTPresponse{
			Status: status.SuccessBstatus,
		})
	}
}

func (h handler) UpdatePerson() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.UpdatePersonParams
		if err := ctx.BodyParser(&params); err != nil {
			return status.ErrInput
		}
		err := h.uc.UpdatePerson(ctx.Context(), params)
		if err != nil {
			return err
		}
		return ctx.JSON(status.HTTPresponse{
			Status: status.SuccessBstatus,
		})
	}
}

func (h handler) DeletePerson() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		rawPersonId := ctx.Query("person_id")
		personId, err := strconv.ParseInt(rawPersonId, 10, 64)
		if err != nil {
			return status.ErrInput.Wrap(fmt.Errorf("cannot get query parameter - person_id"))
		}
		err = h.uc.DeletePerson(ctx.Context(), personId)
		if err != nil {
			return err
		}
		return ctx.JSON(status.HTTPresponse{
			Status: status.SuccessBstatus,
		})
	}
}

func (h handler) FetchPersons() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.FetchPersonsParams
		if err := ctx.QueryParser(&params); err != nil {
			return status.ErrInput.Wrap(fmt.Errorf("cannot parse input query params"))
		}
		res, err := h.uc.FetchPersons(ctx.Context(), params)
		if err != nil {
			return err
		}
		return ctx.JSON(status.HTTPresponse{
			Status: status.SuccessBstatus,
			Result: res,
		})
	}
}
