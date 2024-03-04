package handlers

import (
	"bwg/internal/models"
	"github.com/gofiber/fiber/v2"
)

type service interface {
	CreateNewTicker(ticker string) error
}

type Handlers struct {
	service
}

func New(service service) *Handlers {
	return &Handlers{service}
}
func (h *Handlers) Post(c *fiber.Ctx) error {
	var NewTicker models.NewTicker
	err := c.BodyParser(&NewTicker)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	err = h.service.CreateNewTicker(NewTicker.Ticker)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
	// Find the note with the given id
	// Return the note with the id
}

func (h *Handlers) Get(c *fiber.Ctx) error {
	return nil
}
