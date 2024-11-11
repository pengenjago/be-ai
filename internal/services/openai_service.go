package services

import (
	"be-ai/config"
	"github.com/sashabaranov/go-openai"
)

var client *openai.Client

func GetOpenAI() *openai.Client {
	if client == nil {
		client = openai.NewClient(config.Get("openai.apikey"))
	}

	return client
}
