package aiquery

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	model *genai.GenerativeModel
}

func NewGeminiClient() (*GeminiClient, error) {

	apiKey := os.Getenv("GEMINI_API_KEYS")

	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY not found")
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-3-flash-preview")

	return &GeminiClient{
		model: model,
	}, nil
}

func (g *GeminiClient) Ask(prompt string) (string, error) {

	ctx := context.Background()

	resp, err := g.model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	if len(resp.Candidates) == 0 {
		return "", fmt.Errorf("No response from Gemini")
	}

	part := resp.Candidates[0].Content.Parts[0]

	return fmt.Sprintf("%v", part), nil
}