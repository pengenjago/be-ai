package handlers

import (
	"be-ai/internal/dto"
	"be-ai/internal/services"
	"be-ai/internal/token"
	"be-ai/util"
	"github.com/gofiber/fiber/v2"
)

var threadHandler *ThreadHandler

type ThreadHandler struct {
}

func GetThreadHandler() *ThreadHandler {
	if threadHandler == nil {
		threadHandler = &ThreadHandler{}
	}
	return threadHandler
}

func ThreadRoutes(r fiber.Router, h *ThreadHandler) {
	r.Get("/threads", token.Allow(), h.getAll)
	r.Post("/threads", token.Allow(), h.create)
}

// ------------------------------------------
func (h *ThreadHandler) create(c *fiber.Ctx) error {
	auth := token.GetInfoAuth(c)

	var req dto.ThreadReq
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, err)
	}

	threadId, err := services.GetThreadService().Create(auth.ID, req.AssistantId)
	if err != nil {
		return util.SendError(c, err)
	}

	return util.SendResult(c, threadId)
}

func (h *ThreadHandler) getAll(c *fiber.Ctx) error {
	auth := token.GetInfoAuth(c)

	data := services.GetThreadService().GetAll(auth.ID)
	if data == nil {
		return util.SendNotFound(c)
	}

	return util.SendResult(c, data)
}
