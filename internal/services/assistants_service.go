package services

import (
	"be-ai/internal/constants"
	"be-ai/internal/dto"
	"be-ai/internal/models"
	"be-ai/internal/repositories"
	"context"
	"github.com/sashabaranov/go-openai"
	"log"
)

type AssistantsService interface {
	GetAll() []dto.AssistantsRes
	Create(param dto.AssistantsReq) error
	UploadFile(uploadReq dto.UploadReq) (string, error)
}

var as *assistantsServiceImpl

type assistantsServiceImpl struct {
	assistRepo repositories.AssistantsRepository
}

func GetAssistantsService() AssistantsService {
	if as == nil {
		as = &assistantsServiceImpl{
			assistRepo: repositories.GetAssistantsRepo(),
		}
	}
	return as
}

// ------------------------------------------

func (s *assistantsServiceImpl) GetAll() []dto.AssistantsRes {
	var res []dto.AssistantsRes

	data := s.assistRepo.GetAll()
	for _, val := range data {
		res = append(res, dto.AssistantsRes{
			ID:           val.ID,
			Name:         val.Name,
			Instructions: val.Instructions,
			GptModel:     val.GptModel,
			VectorID:     val.VectorID,
		})
	}

	return res
}

func (s *assistantsServiceImpl) Create(param dto.AssistantsReq) error {
	ai := GetOpenAI()

	resp, err := ai.CreateAssistant(context.Background(), openai.AssistantRequest{
		Name:         &param.Name,
		Instructions: &param.Instructions,
		Model:        openai.GPT4oMini,
		Tools: []openai.AssistantTool{
			{
				Type: openai.AssistantToolTypeFileSearch,
			},
		},
	})
	if err != nil {
		log.Println("failed create assistants openai :", err.Error())
		return constants.ErrConnectOpenAI
	}

	vector, err := ai.CreateVectorStore(context.Background(), openai.VectorStoreRequest{
		Name: param.Name + "_vectorstore",
	})

	assistant := models.Assistants{
		ID:           resp.ID,
		Name:         param.Name,
		Instructions: param.Instructions,
		GptModel:     openai.GPT4oMini,
		VectorID:     vector.ID,
	}

	err = s.assistRepo.Create(&assistant)
	if err != nil {
		log.Println("error create assistants :", err.Error())
		return constants.ErrCreate
	}

	return nil
}

func (s *assistantsServiceImpl) UploadFile(uploadReq dto.UploadReq) (string, error) {
	ai := GetOpenAI()

	file, err := ai.CreateFileBytes(context.Background(), openai.FileBytesRequest{
		Name:    uploadReq.FileName,
		Bytes:   uploadReq.File,
		Purpose: openai.PurposeAssistants,
	})
	if err != nil {
		log.Println("failed create file bytes openai :", err.Error())
		return "", constants.ErrConnectOpenAI
	}

	vectorFile, err := ai.CreateVectorStoreFile(context.Background(), uploadReq.VectorID, openai.VectorStoreFileRequest{
		FileID: file.ID,
	})
	if err != nil {
		log.Println("failed create vector store file openai :", err.Error())
		return "", constants.ErrConnectOpenAI
	}

	_, err = ai.ModifyAssistant(context.Background(), uploadReq.AssistantID, openai.AssistantRequest{
		ToolResources: &openai.AssistantToolResource{
			FileSearch: &openai.AssistantToolFileSearch{
				VectorStoreIDs: []string{uploadReq.VectorID},
			},
		},
	})
	if err != nil {
		return "", constants.ErrConnectOpenAI
	}

	return vectorFile.ID, nil
}
