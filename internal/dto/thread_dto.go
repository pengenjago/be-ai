package dto

import "time"

type ThreadReq struct {
	AssistantId string `json:"assistantId"`
}

type ThreadRes struct {
	ThreadId     string    `json:"threadId"`
	AssistantId  string    `json:"assistantId"`
	CreatedAt    time.Time `json:"createdAt"`
	Topic        string    `json:"topic"`
	Subtopics    []string  `json:"subtopics"`
	Keywords     []string  `json:"keywords"`
	FirstMessage string    `json:"firstMessage"`
	LastMessage  string    `json:"lastMessage"`
	MessageCount int       `json:"messageCount"`
	LastActivity time.Time `json:"lastActivity"`
}
