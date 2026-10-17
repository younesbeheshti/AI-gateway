package ollama

import (
	"context"
	"log"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/yonesbeheshti/ai-gateway/entity"
	"github.com/yonesbeheshti/ai-gateway/pkg/estimatetoken"
)

type Provider struct {
	client *api.Client
}

func New() *Provider {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("Unable to create client: %s", err.Error())
	}
	return &Provider{client: client}
}

func (p *Provider) Client() *api.Client {
	return p.client
}

func (p *Provider) Chat(ctx context.Context, req *entity.ChatRequest) (*entity.ChatResponse, error) {

	messages := make([]api.Message, 0, len(req.Messages))
	for _, msg := range req.Messages {
		messages = append(messages, api.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	log.Println("msgs", messages)

	// TODO: add to config
	chReq := &api.ChatRequest{
		Model:    "gemma2:2b",
		Messages: messages,
	}

	var sb strings.Builder
	var promptTokens, completionTokens int64
	respFunc := func(resp api.ChatResponse) error {

		if resp.PromptEvalCount > 0 {
			promptTokens = int64(resp.PromptEvalCount)
		}

		if resp.EvalCount > 0 {
			completionTokens = int64(resp.EvalCount)
		}

		sb.WriteString(resp.Message.Content)
		return nil
	}

	err := p.client.Chat(ctx, chReq, respFunc)
	if err != nil {
		log.Printf("Unable to send message: %s\n", err.Error())
		return nil, err
	}

	if promptTokens == 0 {
		for _, m := range req.Messages {
			promptTokens += estimatetoken.EstimateTokens(m.Content)
		}
	}
	if completionTokens == 0 {
		completionTokens = estimatetoken.EstimateTokens(sb.String())
	}

	return &entity.ChatResponse{
		Content: sb.String(),
		Usage: entity.Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
	}, nil
}
