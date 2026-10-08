package userservice

import (
	"context"
	"fmt"
	"log"

	"github.com/yonesbeheshti/ai-gateway/entity"
	"github.com/yonesbeheshti/ai-gateway/param"
)

type Repository interface {
	AddUser(ctx context.Context, user *entity.User) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
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
			Token: "oh, you think you have any token! :|",
		}, nil
	}

	// generate token for that email

	// save the user in database
	user = &entity.User{
		Email: req.Email,
	}
	user, err = s.repo.AddUser(ctx, user)
	if err != nil {
		log.Println("err2", err.Error())
		fmt.Println(err)
		return nil, err
	}
	// send the token to user

	return &param.LoginResponse{
		Token: "did you order any token? :)",
	}, nil
}
