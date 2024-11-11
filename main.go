package main

import (
	"be-ai/config"
	"be-ai/server"
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	s := server.NewServer()
	s.App = fiber.New(fiber.Config{
		AppName: config.Get("app.name"),
	})

	db, err := config.ConnectDb()
	if err != nil {
		log.Println("error connect to db :", err.Error())
		return
	}
	s.Db = db

	s.Start()
}
