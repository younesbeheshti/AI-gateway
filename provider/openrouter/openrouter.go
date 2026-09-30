package openrouter

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	openrouter "github.com/OpenRouterTeam/go-sdk"
	"github.com/OpenRouterTeam/go-sdk/models/components"
	"github.com/yonesbeheshti/ai-gateway/entity"
)

type Provider struct {
	client *openrouter.OpenRouter
}

func New() *Provider {
	// TODO: add to config
	//_ = godotenv.Load(".env")
	fmt.Println(os.Getenv("OPENROUTER_API_KEY"))
	client := openrouter.New(
		openrouter.WithSecurity(os.Getenv("OPENROUTER_API_KEY")),
	)
	return &Provider{client: client}
}

func (p *Provider) Client() *openrouter.OpenRouter {
	return p.client
}

func (p *Provider) Chat(ctx context.Context, req *entity.ChatRequest) (*entity.ChatResponse, error) {

	t1 := time.Now()

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

	//fmt.Println(req)
	//// TODO: add to config
	//res, err := p.client.Chat.Send(ctx, components.ChatRequest{
	//	Model: openrouter.Pointer("nvidia/nemotron-3-super-120b-a12b:free"),
	//	Messages: []components.ChatMessages{
	//		components.CreateChatMessagesUser(
	//			components.ChatUserMessage{
	//				Role: components.ChatUserMessageRoleUser,
	//				//Role: "user",
	//				Content: components.CreateChatUserMessageContentStr(
	//					"what is the most dangerous animal on the Earth",
	//				),
	//			},
	//		),
	//	},
	//	Stream: openrouter.Pointer(true),
	//}, nil)
	if err != nil {
		log.Println(err)
	}
	if res == nil || res.EventStream == nil {
		log.Println("expected streaming response")
	}

	stream := res.EventStream
	defer stream.Close()

	var sb strings.Builder

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
	fmt.Println()

	if err := stream.Err(); err != nil {
		log.Println(err)
	}

	fmt.Println(sb.String())
	fmt.Println("time to response:", time.Since(t1))
	return &entity.ChatResponse{
		Content: sb.String(),
	}, nil
}

//func main() {
//	a := New()
//	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
//	defer cancel()
//	req, _ := a.Chat(ctx, nil)
//	fmt.Println(req)
//}
