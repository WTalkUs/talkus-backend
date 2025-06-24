package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"strings"
)

type TextSimilarityRequest struct {
	PostText  string `json:"post_text"`
	GroupText string `json:"group_text"`
}

type TextSimilarityResponse struct {
	Similarity float64 `json:"similarity"`
	Error      string  `json:"error,omitempty"`
}

// CalculateTextSimilarity calculates the semantic similarity between the post text and the group description
func CalculateTextSimilarity(postText, groupText string) (float64, string, error) {
	// Normalizar textos
	postText = strings.TrimSpace(postText)
	groupText = strings.TrimSpace(groupText)

	if postText == "" || groupText == "" {
		return 0.0, "", fmt.Errorf("texts cannot be empty")
	}

	// Obtener embeddings usando Hugging Face
	embeddings, err := getEmbeddings([]string{groupText, postText})
	if err != nil {
		return 0.0, "", fmt.Errorf("error getting embeddings: %v", err)
	}

	if len(embeddings) < 2 {
		return 0.0, "", fmt.Errorf("not enough embeddings obtained")
	}

	// Calcular similitud coseno
	similarity := cosineSimilarity(embeddings[0], embeddings[1])
	veredict := GetContentVerdict(similarity)

	return similarity, veredict, nil
}

// getEmbeddings gets embeddings of texts using the Hugging Face API
func getEmbeddings(texts []string) ([][]float64, error) {
	// Preparar el request
	bodyData := map[string][]string{
		"inputs": texts,
	}
	bodyJSON, err := json.Marshal(bodyData)
	if err != nil {
		return nil, fmt.Errorf("error serializing request: %v", err)
	}

	// Crear request HTTP
	req, err := http.NewRequest("POST",
		"https://router.huggingface.co/hf-inference/models/sentence-transformers/all-MiniLM-L6-v2/pipeline/feature-extraction",
		bytes.NewBuffer(bodyJSON))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("HUGGINGFACE_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	// Ejecutar request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	// Leer respuesta
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Verificar status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: %s - %s", resp.Status, string(bodyBytes))
	}

	// Parsear embeddings
	var embeddings [][]float64
	if err := json.Unmarshal(bodyBytes, &embeddings); err != nil {
		return nil, fmt.Errorf("error parsing embeddings: %v - %s", err, string(bodyBytes))
	}

	return embeddings, nil
}

// cosineSimilarity calculates the cosine similarity between two vectors
func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0.0
	}

	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	denominator := math.Sqrt(normA) * math.Sqrt(normB)
	if denominator == 0 {
		return 0.0
	}

	return dot / denominator
}

// GetContentVerdict returns a verdict based on the similarity score
func GetContentVerdict(similarity float64) string {
	if similarity >= 0.5 {
		return "✅ The post is relevant to the group"
	} else if similarity >= 0.35 {
		return "⚠️ Could be related, requires review"
	} else {
		return "❌ The post is not related to the group"
	}
}
