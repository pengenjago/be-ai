package util

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

const (
	SUCCESS string = "Success"
)

func SendPaged(ctx *fiber.Ctx, data interface{}, pageNo int, pageSize int, totalRecord int) error {
	totalPage := 0

	if pageSize > 0 {
		totalPage = totalRecord / pageSize
		if totalRecord%pageSize > 0 {
			totalPage++
		}
	}

	if pageNo > totalPage {
		pageNo = totalPage
	}

	return ctx.Status(fiber.StatusOK).JSON(ResponsePaged{Status: fiber.StatusOK, Message: SUCCESS, PageNo: pageNo, PageSize: pageSize, TotalRecord: totalRecord, PageTotal: totalPage, Data: data})
}

func SendResult(ctx *fiber.Ctx, data interface{}) error {
	return ctx.Status(fiber.StatusOK).JSON(ResponseData{Status: fiber.StatusOK, Message: SUCCESS, Data: data})
}

func SendError(ctx *fiber.Ctx, err error) error {
	if strings.Contains(err.Error(), "not found") {
		return SendNotFound(ctx)
	}
	return ctx.Status(fiber.StatusBadRequest).JSON(Response{Status: fiber.StatusBadRequest, Message: err.Error()})
}

func SendNotFound(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotFound).JSON(ResponseData{Status: fiber.StatusNotFound, Message: "data not found"})
}

func SendUnauth(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(Response{Status: fiber.StatusUnauthorized, Message: "access not allowed"})
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type ResponseData struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponsePaged struct {
	Status      int         `json:"status"`
	Message     string      `json:"message"`
	Data        interface{} `json:"data"`
	PageNo      int         `json:"pageNo"`
	PageSize    int         `json:"pageSize"`
	PageTotal   int         `json:"pageTotal"`
	TotalRecord int         `json:"totalRecord"`
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(hash), err
}

func CheckPassword(password, hasPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hasPassword), []byte(password))
}
