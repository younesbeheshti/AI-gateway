package aiservice

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yonesbeheshti/ai-gateway/entity"
	"github.com/yonesbeheshti/ai-gateway/param"
	"github.com/yonesbeheshti/ai-gateway/pkg/cache"
)

type Cache interface {
	Get(ctx context.Context, key string) (*entity.ChatResponse, bool, error)
	Set(ctx context.Context, key string, resp *entity.ChatResponse, ttl time.Duration) error
}

type Provider interface {
	Chat(ctx context.Context, req *entity.ChatRequest) (*entity.ChatResponse, error)
}

type Config struct {
	CacheTtl time.Duration `koanf:"cache_ttl"`
}

type Service struct {
	providers       map[string]Provider
	defaultProvider Provider
	cache           Cache
	cfg             Config
}

func New(cfg Config, defaultProvider Provider, providers map[string]Provider, cache Cache) Service {
	return Service{
		providers:       providers,
		defaultProvider: defaultProvider,
		cache:           cache,
		cfg:             cfg,
	}

}

func (s *Service) GenerateResponse(ctx context.Context, req *param.ChatRequest) (*param.ChatResponse, error) {

	// set provider or use default provider
	provider := s.resolveProvider(req.Model)

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

	// check redis to see if the user request or similar to it exist
	var cacheKey string
	if s.cache != nil && req.WithCache {
		if key, err := cache.GenerateCacheKey(userReq); err == nil {
			cacheKey = key
			if cachedResp, hit, err := s.cache.Get(ctx, cacheKey); err == nil && hit {
				log.Println("Cache HIT for key:", cacheKey)
				return &param.ChatResponse{Content: cachedResp.Content}, nil
			}
		}
	}

	// if not, send the request to the selected provider
	// if the selected provider did not respond, try two more times and then change the provider or select default provider
	// 1. Try selected provider
	resp, err := s.chatWithRetry(ctx, provider, userReq, 2)

	// 2. Fall back to default provider only if different and primary failed
	if err != nil && provider != s.defaultProvider && s.defaultProvider != nil {
		log.Printf("Primary provider failed (%v). Falling back to default provider...", err)
		resp, err = s.chatWithRetry(ctx, s.defaultProvider, userReq, 1)
	}

	// 3. Check error before accessing resp
	if err != nil {
		return nil, fmt.Errorf("failed to generate response: %w", err)
	}

	if s.cache != nil && cacheKey != "" {
		if err := s.cache.Set(ctx, cacheKey, resp, s.cfg.CacheTtl); err != nil {
			log.Printf("Warning: failed to save to cache: %v\n", err)
		}
		log.Println("Cached:", cacheKey)
	}
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

func (s *Service) chatWithRetry(ctx context.Context, p Provider, req *entity.ChatRequest, maxRetries int) (*entity.ChatResponse, error) {
	var resp *entity.ChatResponse
	var err error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Stop if client canceled or timed out
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		resp, err = p.Chat(ctx, req)
		if err == nil && resp != nil {
			return resp, nil
		}

		log.Printf("Provider attempt %d failed: %v", attempt+1, err)

		// Wait briefly before retrying (exponential backoff)
		if attempt < maxRetries {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(100*(1<<attempt)) * time.Millisecond):
			}
		}
	}

	return nil, err
}
