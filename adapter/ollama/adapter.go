package main

import (
	"context"
	"log"
	"strings"

	"github.com/ollama/ollama/api"
	"github.com/yonesbeheshti/ai-gateway/entity"
)

type Adapter struct {
	client *api.Client
}

func New() *Adapter {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatalf("Unable to create client: %s", err.Error())
	}
	return &Adapter{client: client}
}

func (a *Adapter) Chat(ctx context.Context, req *entity.ChatRequest) (*entity.ChatResponse, error) {

	messages := make([]api.Message, 0, len(req.Messages))
	for _, msg := range req.Messages {
		messages = append(messages, api.Message{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// TODO: add to config
	chReq := &api.ChatRequest{
		Model:    "gemma2:2b",
		Messages: messages,
	}

	var sb strings.Builder
	respFunc := func(resp api.ChatResponse) error {
		sb.WriteString(resp.Message.Content)
		return nil
	}

	err := a.client.Chat(ctx, chReq, respFunc)
	if err != nil {
		log.Fatalf("Unable to send message: %s", err.Error())
		return nil, err
	}

	return &entity.ChatResponse{
		Content: sb.String(),
	}, nil
}

//func main() {
//	client := New()
//	resp, _ := client.Chat(context.Background(), &entity.ChatRequest{
//		Model: "gemma2:2b",
//		Messages: []entity.Message{
//			{Role: "user", Content: "print all the first ten prime numbers"},
//		},
//		Stream: false,
//	})
//
//	fmt.Println(resp.Content)
//}
