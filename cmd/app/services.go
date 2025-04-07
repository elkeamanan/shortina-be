package app

import "elkeamanan/shortina/internal/link/service"

type ServiceContainer struct {
	LinkService service.LinkService
}

func InitServices(repositories *RepositoryContainer) (ServiceContainer, error) {
	linkService := service.NewLinkService(repositories.LinkRepostitory)
	return ServiceContainer{
		LinkService: linkService,
	}, nil
}
