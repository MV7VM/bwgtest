package handlers

import (
	"bwg/internal/models"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

type service interface {
	CreateNewTicker(ticker string) error
	GetTickerInfo(info models.TickerInfo) (models.TicketDifference, error)
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
	var Params models.TickerInfo
	err := c.BodyParser(&Params)
	slog.Info("Get params: ", Params)
	if err != nil {
		slog.Error("Error in body parser:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	Response, err := h.service.GetTickerInfo(Params)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	response, _ := json.Marshal(Response)
	return c.SendString(string(response))
}
