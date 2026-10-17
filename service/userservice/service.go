package userservice

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/yonesbeheshti/ai-gateway/entity"
	"github.com/yonesbeheshti/ai-gateway/param"
)

type Repository interface {
	AddUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByApiKey(ctx context.Context, apiKey string) (*entity.User, error)
	DeductUserTokens(ctx context.Context, userID uint64, tokens int64) (int64, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}

}

func (s *Service) Login(ctx context.Context, req *param.LoginRequest) (*param.LoginResponse, error) {

	//validate user request
	// check to see if the user exist

	fmt.Println("here")
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Println("err1", err.Error())
		fmt.Println(err)
		return nil, err
	} else if user != nil {
		return &param.LoginResponse{
			ApiKey: user.ApiKey,
			Token:  user.TokenBalance,
		}, nil
	}

	// generate token for that email
	apikey, err := GenerateApiKey()
	if err != nil {
		log.Println("err:", err.Error())
		return nil, err
	}
	// save the user in database
	user = &entity.User{
		Email:        req.Email,
		ApiKey:       apikey,
		TokenBalance: 100_000,
	}
	user, err = s.repo.AddUser(ctx, user)
	if err != nil {
		log.Println("err2", err.Error())
		fmt.Println(err)
		return nil, err
	}
	// send the token to user

	return &param.LoginResponse{
		ApiKey: user.ApiKey,
		Token:  user.TokenBalance,
	}, nil
}

func (s *Service) GetUserByApiKey(ctx context.Context, apiKey string) (*param.User, error) {
	user, err := s.repo.GetUserByApiKey(ctx, apiKey)
	if err != nil {
		log.Println("err", err.Error())
		return nil, err
	}

	userResp := &param.User{
		ID:     user.ID,
		Tokens: user.TokenBalance,
	}

	return userResp, nil
}

func (s *Service) DeductTokens(ctx context.Context, userID uint64, tokens int64) (int64, error) {
	if tokens <= 0 {
		return 0, nil
	}
	return s.repo.DeductUserTokens(ctx, userID, tokens)
}

func GenerateApiKey() (string, error) {
	bytes := make([]byte, 24) // 24 bytes = 48 hex characters
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return "sk-gw-" + hex.EncodeToString(bytes), nil
}
