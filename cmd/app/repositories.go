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
		Host:     config.Cfg.Database.Host,
		Port:     config.Cfg.Database.Port,
		DBName:   config.Cfg.Database.DBName,
		Username: config.Cfg.Database.Username,
		Password: config.Cfg.Database.Password,
		ConnectionConfig: postgres.ConnectionConfig{
			PingTimeout:  config.Cfg.Database.Connection.PingTimeout,
			MaxIdleTime:  config.Cfg.Database.Connection.MaxIdleTime,
			MaxLifetime:  config.Cfg.Database.Connection.MaxLifetime,
			MaxOpenConns: config.Cfg.Database.Connection.MaxOpenConns,
			MaxIdleConns: config.Cfg.Database.Connection.MaxIdleConns,
		},
		MigrationConfig: postgres.MigrationConfig{
			RunMigration:   config.Cfg.Database.Migration.RunMigration,
			MigrationsPath: config.Cfg.Database.Migration.Path,
		},
		SSLConfig: postgres.SSLConfig{
			SSLMode:     config.Cfg.Database.SSL.SSLMode,
			SSLRootCert: config.Cfg.Database.SSL.SSLRootCert,
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
