package app

import (
	"context"
	"elkeamanan/shortina/config"
	linkRepository "elkeamanan/shortina/internal/link/repository"
	userRepository "elkeamanan/shortina/internal/user/repository"
	"elkeamanan/shortina/storage/postgres"
)

type RepositoryContainer struct {
	LinkRepository linkRepository.LinkRepository
	UserRepository userRepository.UserRepository
}

func InitRepositories(ctx context.Context) (*RepositoryContainer, error) {
	st, err := postgres.InitDatabase(ctx, postgres.DatabaseConfig{
		Host:        config.Cfg.Database.Host,
		Port:        config.Cfg.Database.Port,
		DBName:      config.Cfg.Database.DBName,
		Username:    config.Cfg.Database.Username,
		Password:    config.Cfg.Database.Password,
		PingTimeout: config.Cfg.Database.PingTimeout,
		MigrationConfig: postgres.MigrationConfig{
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
