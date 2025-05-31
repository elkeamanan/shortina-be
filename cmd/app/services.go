package app

import (
	linkService "elkeamanan/shortina/internal/link/service"
	userService "elkeamanan/shortina/internal/user/service"
	"elkeamanan/shortina/storage/redis"
)

type ServiceContainer struct {
	LinkService linkService.LinkService
	UserService userService.UserService
}

func InitServices(repositories *RepositoryContainer) (*ServiceContainer, error) {
	redisClient, err := redis.NewRedisClient()
	if err != nil {
		return nil, err
	}

	linkService := linkService.NewLinkService(repositories.LinkRepository)
	userService := userService.NewUserService(repositories.UserRepository, redisClient)
	return &ServiceContainer{
		LinkService: linkService,
		UserService: userService,
	}, nil
}
