package openai

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

const (
	systemPromptPath = "pkg/openai/system_prompt.txt"
)

type OpenAI struct {
	Client *openai.Client
}

type AnalysisResponse struct {
	Security  string // Underlying security to trade
	Sentiment int    // Either -1, 0, or 1
	Summary   string // Small justification for the response to aid in manual validation
}

func NewOpenAI() (*OpenAI, error) {
	client := openai.NewClient(os.Getenv("OpenAIAPIKey"))
	return &OpenAI{
		Client: client,
	}, nil
}

func (o *OpenAI) PerformAnalysisOnPost(post reddit.Post) (*AnalysisResponse, error) {
	userPrompt := redditPostToOpenAIPrompt(post)
	systemPrompt, err := os.ReadFile(systemPromptPath)
	if err != nil {
		return nil, fmt.Errorf("issue reading system prompt file: %v", err)
	}
	response, err := o.Client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    "system",
				Content: string(systemPrompt),
			},
			{
				Role:    "user",
				Content: userPrompt,
			},
		},
		MaxTokens: 100,
	})
	if err != nil {
		return nil, fmt.Errorf("issue getting completion from OpenAI: %w", err)
	}
	analysis, err := stringResponseToAnalysisResponse(response.Choices[0].Message.Content)
	if err != nil {
		return nil, fmt.Errorf("issue converting response from OpenAI into analysis: %w", err)
	}
	return analysis, nil
}

func redditPostToOpenAIPrompt(post reddit.Post) string {
	return strings.Join([]string{post.Title, post.Body}, "\n\n")
}

func stringResponseToAnalysisResponse(response string) (*AnalysisResponse, error) {
	if response == "N/A" {
		fmt.Println("OpenAI returned N/A")
		return nil, nil
	}
	splitResponse := strings.Split(response, ";;") // note that the ;; delimiter is set in the system_prompt.txt file
	if len(splitResponse) != 3 {
		return nil, fmt.Errorf("expected response to be delimited with ;; instead got: %v", response)
	}
	sentiment, err := strconv.Atoi(splitResponse[1])
	if err != nil {
		return nil, fmt.Errorf("sentiment expected to be an integer, got: %v", err)
	}
	return &AnalysisResponse{
		Security:  splitResponse[0],
		Sentiment: sentiment,
		Summary:   splitResponse[2],
	}, nil
}
