package services

import (
	"be-ai/internal/constants"
	"be-ai/internal/dto"
	"be-ai/internal/repositories"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"log"
	"strings"
	"time"
)

type ThreadService interface {
	GetAll(userId string) []dto.ThreadRes
	Create(userId, assistantId string) (string, error)
}

var ts *threadServiceImpl

type threadServiceImpl struct {
	threadRepo repositories.ThreadRepository
}

func GetThreadService() ThreadService {
	if ts == nil {
		ts = &threadServiceImpl{
			threadRepo: repositories.GetThreadRepo(),
		}
	}
	return ts
}

// ------------------------------------------

func (t *threadServiceImpl) Create(userId, assistantId string) (string, error) {
	th, err := GetOpenAI().CreateThread(context.Background(), openai.ThreadRequest{})
	if err != nil {
		log.Println("failed to create thread openai :", err.Error())
		return "", constants.ErrConnectOpenAI
	}

	err = t.threadRepo.Create(userId, th.ID, assistantId)
	if err != nil {
		log.Println("failed to create thread :", err.Error())
		return "", constants.ErrCreate
	}

	return th.ID, nil
}

func (t *threadServiceImpl) GetAll(userId string) []dto.ThreadRes {
	var res []dto.ThreadRes

	data := t.threadRepo.GetAll(userId)
	for _, val := range data {
		summary, _ := t.getThreadSummary(context.Background(), GetOpenAI(), val.ID)
		summary.AssistantId = val.AssistantId
		res = append(res, summary)
	}
	return res
}

func (t *threadServiceImpl) getThreadSummary(ctx context.Context, ai *openai.Client, threadID string) (dto.ThreadRes, error) {
	res := dto.ThreadRes{
		ThreadId: threadID,
	}

	ascending := "asc"

	messages, err := ai.ListMessage(ctx, threadID, nil, &ascending, nil, nil, nil)
	if err != nil {
		return res, fmt.Errorf("failed to list messages: %v", err)
	}

	if len(messages.Messages) == 0 {
		return res, nil
	}

	firstMsg := messages.Messages[0]
	lastMsg := messages.Messages[len(messages.Messages)-1]

	res.CreatedAt = time.Unix(int64(firstMsg.CreatedAt), 0).Local()
	res.MessageCount = len(messages.Messages)
	res.LastActivity = time.Unix(int64(lastMsg.CreatedAt), 0).Local()

	if len(firstMsg.Content) > 0 {
		res.FirstMessage = firstMsg.Content[0].Text.Value
	}
	if len(lastMsg.Content) > 0 {
		res.LastMessage = lastMsg.Content[0].Text.Value
	}

	err2 := t.getTopicCache(ctx, threadID, &res, messages)
	if err2 != nil {
		return res, err2
	}

	return res, nil
}

func (t *threadServiceImpl) getTopicCache(ctx context.Context, threadID string, res *dto.ThreadRes, messages openai.MessagesList) error {
	cacheTopic, isExist := NewTopicCache().Get(threadID)
	if isExist {
		log.Println("cache topic found:", cacheTopic)
		res.Topic = cacheTopic.MainTopic
		res.Subtopics = cacheTopic.Subtopics
		res.Keywords = cacheTopic.Keywords
	} else {
		log.Println("cache topic not found:", threadID)
		analysis, err := analyzeThreadTopic(ctx, GetOpenAI(), messages.Messages)
		if err != nil {
			return fmt.Errorf("failed to analyze topic: %v", err)
		}
		res.Topic = analysis.MainTopic
		res.Subtopics = analysis.Subtopics
		res.Keywords = analysis.Keywords

		NewTopicCache().Set(threadID, analysis)
	}
	return nil
}

func analyzeThreadTopic(ctx context.Context, client *openai.Client, messages []openai.Message) (TopicAnalysis, error) {
	// Siapkan konteks percakapan untuk analisis
	var conversationContext strings.Builder
	for _, msg := range messages {
		if len(msg.Content) > 0 {
			conversationContext.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content[0].Text.Value))
		}
	}

	// Buat prompt untuk menganalisis topik
	systemPrompt := `You are a conversation topic analyzer. Analyze the conversation and:
1. Determine the main topic being discussed (ignore greetings and small talk)
2. Identify 2-3 subtopics if any
3. Extract 3-5 relevant keywords

Return the analysis in this exact JSON format:
{
    "main_topic": "brief topic description",
    "subtopics": ["subtopic1", "subtopic2"],
    "keywords": ["keyword1", "keyword2", "keyword3"]
}
`

	userPrompt := fmt.Sprintf("Analyze this conversation and determine the main topic, subtopics, and keywords:\n\n%s",
		conversationContext.String())

	// Buat request untuk analisis
	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4TurboPreview,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: "json_object",
		},
	})

	if err != nil {
		return TopicAnalysis{}, fmt.Errorf("failed to analyze topic: %v", err)
	}

	log.Println("check :", resp.Choices[0].Message.Content)

	// Parse response JSON
	var analysis TopicAnalysis
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &analysis); err != nil {
		return TopicAnalysis{}, fmt.Errorf("failed to parse topic analysis: %v", err)
	}

	return analysis, nil
}

// -----------------
var topicCache *TopicCache

type TopicAnalysis struct {
	MainTopic string   `json:"main_topic"`
	Subtopics []string `json:"subtopics"`
	Keywords  []string `json:"keywords"`
}

type TopicCache struct {
	Topics     map[string]TopicAnalysis
	TTL        time.Duration
	LastUpdate map[string]time.Time
}

func NewTopicCache() *TopicCache {
	if topicCache == nil {
		topicCache = &TopicCache{
			Topics:     make(map[string]TopicAnalysis),
			TTL:        30 * time.Minute,
			LastUpdate: make(map[string]time.Time),
		}
	}
	return topicCache
}

func (tc *TopicCache) Get(threadID string) (TopicAnalysis, bool) {
	if lastUpdate, exists := tc.LastUpdate[threadID]; exists {
		if time.Since(lastUpdate) < tc.TTL {
			if topic, exists := tc.Topics[threadID]; exists {
				return topic, true
			}
		}
	}
	return TopicAnalysis{}, false
}

func (tc *TopicCache) Set(threadID string, analysis TopicAnalysis) {
	tc.Topics[threadID] = analysis
	tc.LastUpdate[threadID] = time.Now()
}

// -----------------
