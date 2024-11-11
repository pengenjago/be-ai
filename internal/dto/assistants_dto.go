package dto

type AssistantsReq struct {
	Name         string `json:"name,omitempty"`
	Instructions string `json:"instructions,omitempty"`
}

type AssistantsRes struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Instructions string `json:"instructions,omitempty"`
	GptModel     string `json:"gptModel,omitempty"`
	VectorID     string `json:"vectorId,omitempty"`
}

type UploadReq struct {
	File        []byte
	FileName    string
	AssistantID string `form:"assistantId"`
	VectorID    string `form:"vectorId"`
}

type StreamMessage struct {
	Role    string `json:"role"`
	Content any    `json:"content"`
}

type WSResponse struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}
