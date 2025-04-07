package app

type ServiceContainer struct {
}

func InitServices(repositories *RepositoryContainer) (ServiceContainer, error) {
	return ServiceContainer{}, nil
}
