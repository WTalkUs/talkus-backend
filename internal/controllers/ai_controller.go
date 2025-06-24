package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/JuanPidarraga/talkus-backend/internal/usecases"
)

type AIController struct {
	aiUsecase *usecases.AIUsecase
}

func NewAIController(aiUsecase *usecases.AIUsecase) *AIController {
	return &AIController{
		aiUsecase: aiUsecase,
	}
}

type ValidateContentRequest struct {
	Title     string `json:"title"`
	Content   string `json:"content"`
	GroupText string `json:"group_text,omitempty"`
}

type ValidateContentResponse struct {
	Similarity float64 `json:"similarity"`
	Verdict    string  `json:"verdict"`
	Error      string  `json:"error,omitempty"`
}

// ValidateContentHandler handles the content validation using text similarity
func (c *AIController) ValidateContentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req ValidateContentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	if req.Content == "" {
		http.Error(w, "Content is required", http.StatusBadRequest)
		return
	}

	if req.GroupText == "" {
		http.Error(w, "GroupText is required", http.StatusBadRequest)
		return
	}

	similarity, veredict, err := c.aiUsecase.CalculateTextSimilarity(req.Title+" "+req.Content, req.GroupText)

	response := ValidateContentResponse{}
	if err != nil {
		response.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		response.Similarity = similarity
		response.Verdict = veredict
	}

	json.NewEncoder(w).Encode(response)
}
