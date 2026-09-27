package aiservice

import (
	"context"

	"github.com/yonesbeheshti/ai-gateway/param"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) GenerateResponse(ctx context.Context, req *param.ChatRequest) (*param.ChatResponse, error) {

}
