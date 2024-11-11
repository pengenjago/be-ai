package handlers

import (
	"be-ai/internal/dto"
	"be-ai/internal/services"
	"be-ai/util"
	"github.com/gofiber/fiber/v2"
	"log"
)

var userHandler *UserHandler

type UserHandler struct {
}

func GetUserHandler() *UserHandler {
	if userHandler == nil {
		userHandler = &UserHandler{}
	}
	return userHandler
}

func UserRoutes(r fiber.Router, h *UserHandler) {
	r.Post("/users/login", h.login)
	r.Post("/users", h.create)
}

// ------------------------------------------

func (h *UserHandler) create(c *fiber.Ctx) error {
	var req dto.UserReq
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, err)
	}

	err := services.GetUserService().CreateUser(req)
	if err != nil {
		return util.SendError(c, err)
	}

	return util.SendResult(c, req)
}

func (h *UserHandler) login(c *fiber.Ctx) error {
	log.Println("wkakwka")
	var req dto.UserLogin
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, err)
	}

	res, err := services.GetUserService().LoginUser(req)
	if err != nil {
		return util.SendError(c, err)
	}

	return util.SendResult(c, res)
}
