package openrouter

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	openrouter "github.com/OpenRouterTeam/go-sdk"
	"github.com/OpenRouterTeam/go-sdk/models/components"
	"github.com/yonesbeheshti/ai-gateway/entity"
	"github.com/yonesbeheshti/ai-gateway/pkg/estimatetoken"
)

type Config struct {
	APIKey string `koanf:"api_key"`
}

type Provider struct {
	client *openrouter.OpenRouter
}

func New(cfg Config) *Provider {
	// TODO: add to config
	//_ = godotenv.Load(".env")
	//client := openrouter.New(
	//	openrouter.WithSecurity(os.Getenv("OPENROUTER_API_KEY")),
	//)

	client := openrouter.New(
		openrouter.WithSecurity(cfg.APIKey),
	)
	return &Provider{client: client}
}

func (p *Provider) Client() *openrouter.OpenRouter {
	return p.client
}

func (p *Provider) Chat(ctx context.Context, req *entity.ChatRequest) (*entity.ChatResponse, error) {

	// TODO: user can send multiple messages
	messages := make([]components.ChatMessages, 0, len(req.Messages))
	for _, msg := range req.Messages {
		messages = append(
			messages,
			components.CreateChatMessagesUser(
				components.ChatUserMessage{
					Role: components.ChatUserMessageRoleUser,
					Content: components.CreateChatUserMessageContentStr(
						msg.Content,
					),
				},
			),
		)
	}

	res, err := p.client.Chat.Send(ctx, components.ChatRequest{
		Model:    openrouter.Pointer("nvidia/nemotron-3-super-120b-a12b:free"),
		Messages: messages,
		Stream:   openrouter.Pointer(true),
	}, nil)

	if err != nil {
		log.Println(err)
	}
	if res == nil || res.EventStream == nil {
		log.Println("expected streaming response")
		return nil, errors.New("expected streaming response")
	}

	stream := res.EventStream
	defer stream.Close()

	var sb strings.Builder
	var promptTokens, completionTokens int64 = 0, 0

	for stream.Next() {
		chunk := stream.Value()
		if chunk == nil {
			continue
		}
		for _, choice := range chunk.Data.Choices {
			if text, ok := choice.Delta.Content.Get(); ok && text != nil {
				sb.WriteString(*text)
			}
		}
	}

	if err := stream.Err(); err != nil {
		log.Println(err)
	}

	for _, m := range req.Messages {
		promptTokens += estimatetoken.EstimateTokens(m.Content)
	}

	completionTokens = estimatetoken.EstimateTokens(sb.String())

	fmt.Println("promptTokens:", promptTokens)
	fmt.Println("completionTokens:", completionTokens)
	return &entity.ChatResponse{
		Content: sb.String(),
		Usage: entity.Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
	}, nil
}

//func main() {
//	a := New()
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
//	defer cancel()
//	req, _ := a.Chat(ctx, nil)
//	fmt.Println(req)
//}
