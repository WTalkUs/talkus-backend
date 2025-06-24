package usecases

import (
	"github.com/JuanPidarraga/talkus-backend/internal/service"
)

type AIUsecase struct {
}

func NewAIUsecase() *AIUsecase {
	return &AIUsecase{}
}

// CalculateTextSimilarity calculates the similarity between the post text and the group description
func (uc *AIUsecase) CalculateTextSimilarity(postText, groupText string) (float64, string, error) {
	return service.CalculateTextSimilarity(postText, groupText)
}
