package handlers

import (
	"be-ai/internal/dto"
	"be-ai/internal/services"
	"be-ai/internal/token"
	"be-ai/util"
	"context"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/sashabaranov/go-openai"
	"log"
	"time"
)

var assistantHandler *AssistantsHandler

type AssistantsHandler struct {
}

func GetAssistantHandler() *AssistantsHandler {
	if assistantHandler == nil {
		assistantHandler = &AssistantsHandler{}
	}
	return assistantHandler
}

func AssistantRoutes(r fiber.Router, h *AssistantsHandler) {
	r.Post("/assistants", token.Allow(), h.create)
	r.Post("/assistants/upload", token.Allow(), h.upload)
	r.Get("/assistants/chat", websocket.New(h.chatAssistant))
	r.Get("/chat/stream", websocket.New(h.chatStream))
}

// ------------------------------------------

func (h *AssistantsHandler) create(c *fiber.Ctx) error {
	var req dto.AssistantsReq
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, err)
	}

	err := services.GetAssistantsService().Create(req)
	if err != nil {
		return util.SendError(c, err)
	}

	return util.SendResult(c, req)
}

func (h *AssistantsHandler) upload(c *fiber.Ctx) error {
	var req dto.UploadReq
	if err := c.BodyParser(&req); err != nil {
		return util.SendError(c, err)
	}

	fh, err := c.FormFile("file")
	if err != nil {
		return util.SendError(c, err)
	}

	file, err := fh.Open()
	if err != nil {
		return util.SendError(c, err)
	}
	defer file.Close()

	fileBytes := make([]byte, fh.Size)
	_, err = file.Read(fileBytes)
	if err != nil {
		return util.SendError(c, err)
	}

	req.File = fileBytes
	req.FileName = fh.Filename

	res, err := services.GetAssistantsService().UploadFile(req)
	if err != nil {
		return util.SendError(c, err)
	}

	return util.SendResult(c, res)
}

func (h *AssistantsHandler) chatStream(c *websocket.Conn) {
	defer c.Close()
	ai := services.GetOpenAI()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("error read message:", err)
			break
		}
		req := openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: string(msg)},
			},
			Stream: true,
		}

		stream, err := ai.CreateChatCompletionStream(context.Background(), req)
		if err != nil {
			log.Println("error chat stream :", err.Error())
			break
		}
		defer stream.Close()

		for {
			response, err := stream.Recv()
			if err != nil {
				log.Println("Stream receive error:", err)
				break
			}

			if err := c.WriteMessage(websocket.TextMessage, []byte(response.Choices[0].Delta.Content)); err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	}
}

func (h *AssistantsHandler) chatAssistant(c *websocket.Conn) {
	threadId := c.Query("threadId")
	assistantId := c.Query("assistantId")
	ctx := context.Background()

	defer c.Close()

	ai := services.GetOpenAI()

	// load history message first
	if h.loadHistoryMessage(c, ai, ctx, threadId) {
		return
	}

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("error read message:", err)
			break
		}

		_, err = ai.CreateMessage(ctx, threadId, openai.MessageRequest{
			Role:    openai.ChatMessageRoleUser,
			Content: string(msg),
		})
		if err != nil {
			log.Println("error create message openai :", err)
			break
		}

		run, err := ai.CreateRun(ctx, threadId, openai.RunRequest{
			AssistantID: assistantId,
			Model:       openai.GPT4oMini,
		})
		if err != nil {
			log.Println("error create run openai :", err)
			break
		}

		for {
			if h.streamResponseAssistant(c, ai, ctx, threadId, run) {
				break
			}
		}
	}
}

func (h *AssistantsHandler) streamResponseAssistant(c *websocket.Conn, ai *openai.Client, ctx context.Context, threadId string, run openai.Run) bool {
	runStatus, err := ai.RetrieveRun(ctx, threadId, run.ID)
	if err != nil {
		sendError(c, "Failed to retrieve run: "+err.Error())
		return true
	}

	// Sent status to client
	sendWSResponse(c, "status", runStatus.Status)

	if runStatus.Status == "completed" {
		messages, err := ai.ListMessage(ctx, threadId, nil, nil, nil, nil, &run.ID)
		if err != nil {
			sendError(c, "Failed to list messages: "+err.Error())
			return true
		}

		if len(messages.Messages) > 0 {
			sendWSResponse(c, "message", dto.StreamMessage{
				Role:    messages.Messages[0].Role,
				Content: messages.Messages[0].Content[0].Text.Value,
			})
		}
		return true
	}

	time.Sleep(100 * time.Millisecond)
	return false
}

func (h *AssistantsHandler) loadHistoryMessage(c *websocket.Conn, ai *openai.Client, ctx context.Context, threadId string) bool {

	// List history msg
	messages, err := ai.ListMessage(ctx, threadId, nil, nil, nil, nil, nil)
	if err != nil {
		sendError(c, "Failed to retrieve messages: "+err.Error())
		return true
	}

	// Sent history msg to client
	for _, msg := range messages.Messages {
		sendWSResponse(c, "history", dto.StreamMessage{
			Role:    msg.Role,
			Content: msg.Content[0].Text.Value,
		})
	}
	return false
}

func sendWSResponse(c *websocket.Conn, msgType string, message interface{}) {
	response := dto.WSResponse{
		Type:    msgType,
		Message: message,
	}

	if err := c.WriteJSON(response); err != nil {
		log.Printf("Error sending WebSocket message: %v", err)
	}
}

func sendError(c *websocket.Conn, message string) {
	sendWSResponse(c, "error", message)
}
