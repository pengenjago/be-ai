package server

import (
	"be-ai/config"
	"be-ai/internal/handlers"
	"be-ai/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/gorm"
	"time"
)

type Server struct {
	App *fiber.App
	Db  *gorm.DB
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Start() {

	s.App.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Success running at :" + time.Now().Format("2006-01-02 15:04:05")})
	})

	s.App.Use(cors.New(cors.Config{
		AllowHeaders: "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	s.migrate()
	s.handlers()

	_ = s.App.Listen(":" + config.Get("app.port"))
}

func (s *Server) migrate() {
	_ = s.Db.AutoMigrate(models.Assistants{})
	_ = s.Db.AutoMigrate(models.User{})
	_ = s.Db.AutoMigrate(models.Thread{})
}

func (s *Server) handlers() {
	api := s.App.Group("/api")

	handlers.AssistantRoutes(api, handlers.GetAssistantHandler())
	handlers.UserRoutes(api, handlers.GetUserHandler())
	handlers.ThreadRoutes(api, handlers.GetThreadHandler())
}
