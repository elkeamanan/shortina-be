package service

import (
	"context"
	"elkeamanan/shortina/config"
	"elkeamanan/shortina/internal/user/domain"
	"elkeamanan/shortina/internal/user/repository"
	"elkeamanan/shortina/storage/redis"
	"fmt"

	"github.com/google/uuid"
)

type userService struct {
	userRepo    repository.UserRepository
	redisClient redis.RedisClient
}

func NewUserService(userRepo repository.UserRepository, redisClient redis.RedisClient) UserService {
	return &userService{userRepo: userRepo, redisClient: redisClient}
}

func (s *userService) RegisterUser(ctx context.Context, req *domain.RegisterUserRequest) error {
	err := req.Validate()
	if err != nil {
		return err
	}

	newUser, err := req.ToNewUser()
	if err != nil {
		return err
	}

	return s.userRepo.CreateUser(ctx, nil, newUser)
}

func (s *userService) LoginUser(ctx context.Context, req *domain.LoginUserRequest) (*domain.TokenPair, error) {
	err := req.Validate()
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetUser(ctx, domain.UserPredicate{Email: req.Email})
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found with email %s", req.Email)
	}

	if !domain.ValidatePassword([]byte(user.Password), []byte(req.Password)) {
		return nil, fmt.Errorf("invalid password")
	}

	tokenPair, err := domain.GenerateTokenPair(user)
	if err != nil {
		return nil, err
	}

	err = s.redisClient.Set(ctx, fmt.Sprintf("refresh_token:%s", tokenPair.RefreshToken), user.ID.String(), config.Cfg.Token.RefreshTokenExpiry)
	if err != nil {
		return nil, err
	}

	return tokenPair, err
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, req *domain.UpdateUserRequest) error {
	pred := domain.UserPredicate{ID: id}
	requestedUser, err := s.userRepo.GetUser(ctx, pred)
	if err != nil {
		return err
	}

	if requestedUser == nil {
		return fmt.Errorf("user not found with id %s", id.String())
	}

	return s.userRepo.UpdateUser(ctx, nil, &domain.User{Fullname: req.Fullname}, pred)
}

func (s *userService) RefreshUserToken(ctx context.Context, refreshToken string) (*domain.TokenPair, error) {
	userID, err := s.redisClient.Get(ctx, fmt.Sprintf("refresh_token:%s", refreshToken))
	if err != nil {
		return nil, err
	}

	if userID == "" {
		return nil, fmt.Errorf("invalid refresh token")
	}
	user, err := s.userRepo.GetUser(ctx, domain.UserPredicate{ID: uuid.MustParse(userID)})
	if err != nil {
		return nil, err
	}

	return domain.GenerateTokenPair(user)
}
