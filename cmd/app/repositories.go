package app

import (
	"context"
	"elkeamanan/shortina/cmd/storage"
	"elkeamanan/shortina/config"
	linkRepository "elkeamanan/shortina/internal/link/repository"
	userRepository "elkeamanan/shortina/internal/user/repository"
)

type RepositoryContainer struct {
	LinkRepository linkRepository.LinkRepository
	UserRepository userRepository.UserRepository
}

func InitRepositories(ctx context.Context) (*RepositoryContainer, error) {
	st, err := storage.InitDatabase(ctx, storage.DatabaseConfig{
		Host:        config.Cfg.Database.Host,
		Port:        config.Cfg.Database.Port,
		DBName:      config.Cfg.Database.DBName,
		Username:    config.Cfg.Database.Username,
		Password:    config.Cfg.Database.Password,
		PingTimeout: config.Cfg.Database.PingTimeout,
		MigrationConfig: storage.MigrationConfig{
			RunMigration:   config.Cfg.Database.Migration.RunMigration,
			MigrationsPath: config.Cfg.Database.Migration.Path,
		},
	})
	if err != nil {
		return nil, err
	}

	linkRepo := linkRepository.NewLinkRepository(st)
	userRepo := userRepository.NewUserRepository(st)
	return &RepositoryContainer{
		LinkRepository: linkRepo,
		UserRepository: userRepo,
	}, nil
}
