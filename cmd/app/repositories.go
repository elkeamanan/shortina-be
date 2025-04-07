package app

import (
	"context"
	"elkeamanan/shortina/cmd/storage"
	"elkeamanan/shortina/config"
)

type RepositoryContainer struct {
}

func InitRepositories(ctx context.Context) (*RepositoryContainer, error) {
	_, err := storage.InitDatabase(ctx, storage.DatabaseConfig{
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

	return &RepositoryContainer{}, nil
}
