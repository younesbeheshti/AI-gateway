package aiservice

import (
	"context"

	"github.com/yonesbeheshti/ai-gateway/entity"
	"github.com/yonesbeheshti/ai-gateway/param"
)

type Provider interface {
	Chat(ctx context.Context, req *entity.ChatRequest) (*entity.ChatResponse, error)
}

type Service struct {
	providers       map[string]Provider
	defaultProvider Provider
}

func New(defaultProvider Provider, providers map[string]Provider) Service {
	return Service{
		providers:       providers,
		defaultProvider: defaultProvider,
	}
}

func (s *Service) GenerateResponse(ctx context.Context, req *param.ChatRequest) (*param.ChatResponse, error) {

	// set provider or use default provider
	provider := s.resolveProvider(req.Model)

	// check redis to see if the user request or similar to it exist
	// if not, send the request to the selected provider
	messages := make([]entity.Message, 0)
	for _, message := range req.Messages {
		msg := entity.Message{
			Role:    message.Role,
			Content: message.Content,
		}
		messages = append(messages, msg)
	}
	userReq := &entity.ChatRequest{
		Model:    req.Model,
		Messages: messages,
	}
	resp, err := provider.Chat(ctx, userReq)
	if err != nil {
		return nil, err
	}
	// if the selected provider did not respond, try two more times and then change the provider or select default provider
	return &param.ChatResponse{
		Content: resp.Content,
	}, nil
}

func (s *Service) resolveProvider(model string) Provider {
	// check if the model is available
	if p, ok := s.providers[model]; ok && p != nil {
		return s.providers[model]
	}
	return s.defaultProvider
}
