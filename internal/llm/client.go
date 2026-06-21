package llm

import (
	"context"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"
)

type LLM struct {
	Model                string
	APIKey               string
	BaseURL              string
	Temperature          float64
	MaxTokens            int
	TotalPromptTokens    int
	TotalCompletionTokens int
	client               *openai.Client
}

func NewLLM(model, apiKey, baseURL string, temperature float64, maxTokens int) *LLM {
	cfg := openai.DefaultConfig(apiKey)
	if baseURL != "" {
		cfg.BaseURL = baseURL
	}
	client := openai.NewClientWithConfig(cfg)
	return &LLM{
		Model:       model,
		APIKey:      apiKey,
		BaseURL:     baseURL,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		client:      client,
	}
}

func (l *LLM) Chat(messages []map[string]interface{}, tools []map[string]interface{}, onToken func(string)) (*LLMResponse, error) {
	openaiMessages := make([]openai.ChatCompletionMessage, 0, len(messages))
	for _, m := range messages {
		role, _ := m["role"].(string)
		content, _ := m["content"].(string)
		openaiMessages = append(openaiMessages, openai.ChatCompletionMessage{
			Role:    role,
			Content: content,
		})
	}

	req := openai.ChatCompletionRequest{
		Model: l.Model,
		Messages: openaiMessages,
		MaxTokens: l.MaxTokens,
	}

	if l.Temperature > 0 {
		req.Temperature = float32(l.Temperature)
	}

	if len(tools) > 0 {
		req.Tools = make([]openai.Tool, len(tools))
		for i, t := range tools {
			fn, _ := t["function"].(map[string]interface{})
			req.Tools[i] = openai.Tool{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name:        fn["name"].(string),
					Description: fn["description"].(string),
					Parameters:  fn["parameters"].(map[string]interface{}),
				},
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Use non-streaming for compatibility
	resp, err := l.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("chat completion error: %w", err)
	}

	content := ""
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
		content = resp.Choices[0].Message.Content
		if onToken != nil {
			onToken(content)
		}
	}

	result := &LLMResponse{
		Content: content,
	}

	if resp.Usage.PromptTokens > 0 {
		result.PromptTokens = resp.Usage.PromptTokens
		result.CompletionTokens = resp.Usage.CompletionTokens
		l.TotalPromptTokens += resp.Usage.PromptTokens
		l.TotalCompletionTokens += resp.Usage.CompletionTokens
	}

	return result, nil
}

func (l *LLM) EstimatedCost() *float64 {
	p, ok := pricing[l.Model]
	if !ok {
		return nil
	}
	cost := (float64(l.TotalPromptTokens) * p[0] / 1_000_000) +
		(float64(l.TotalCompletionTokens) * p[1] / 1_000_000)
	return &cost
}
