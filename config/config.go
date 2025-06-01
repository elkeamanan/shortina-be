package config

import (
	"context"
	"elkeamanan/shortina/provider/secret"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var Cfg Config

func LoadConfig(ctx context.Context) error {
	if err := godotenv.Load(); err != nil {
		log.Error(err.Error())
	}

	err := env.Parse(&Cfg)
	if err != nil {
		return err
	}

	if !Cfg.IsEnableLoadingSecret {
		return nil
	}

	manager, err := secret.NewAWSSecretsManager(ctx, Cfg.Region)
	if err != nil {
		return err
	}

	err = secret.LoadValueFromSecret(ctx, manager, &Cfg)
	if err != nil {
		return err
	}

	return nil

}

type Config struct {
	Environment string `env:"ENVIRONMENT,required"`
	Region      string `env:"SERVICE_REGION,required"`
	ServiceName string `env:"SERVICE_NAME" envDefault:""`
	Timezone    string `env:"TIMEZONE,required" envDefault:"Asia/Jakarta"`

	HttpPort int    `env:"API_PORT" envDefault:"8080"`
	Host     string `env:"GRPC_HOST" envDefault:"localhost"`

	IsEnableLoadingSecret bool `env:"ENABLE_LOADING_SECRET" envDefault:"false"`

	Database struct {
		Host       string `env:"DB_HOST" secret:"SHORTINA_DB_HOST"`
		Port       int    `env:"DB_PORT" envDefault:"5432"`
		DBName     string `env:"DB_NAME" envDefault:"shortina"`
		Username   string `env:"DB_USERNAME" secret:"SHORTINA_DB_USERNAME"`
		Password   string `env:"DB_PASSWORD" secret:"SHORTINA_DB_PASSWORD"`
		Connection struct {
			PingTimeout  time.Duration `env:"DB_TIMEOUT" envDefault:"5s"`
			MaxIdleTime  time.Duration `env:"DB_MAX_IDLE_TIME" envDefault:"2m"`
			MaxLifetime  time.Duration `env:"DB_MAX_LIFETIME" envDefault:"1h"`
			MaxOpenConns int           `env:"DB_MAX_OPEN_CONNS" envDefault:"10"`
			MaxIdleConns int           `env:"DB_MAX_IDLE_CONNS" envDefault:"5"`
		}
		SSL struct {
			SSLMode     string `env:"DB_SSL_MODE" envDefault:"disable"`
			SSLRootCert string `env:"DB_SSL_ROOT_CERT"`
		}
		Migration struct {
			RunMigration bool   `env:"DB_RUN_MIGRATION" envDefault:"true"`
			Path         string `env:"DB_MIGRATION_PATH" envDefault:"storage/postgres/migrations"`
		}
	}

	Redis struct {
		Host string `env:"REDIS_HOST" secret:"SHORTINA_REDIS_HOST"`
		Port int    `env:"REDIS_PORT" envDefault:"6379"`
	}

	Token struct {
		SecretKey          string        `env:"TOKEN_SECRET_KEY" secret:"SHORTINA_TOKEN_SECRET_KEY"`
		AccessTokenExpiry  time.Duration `env:"ACCESS_TOKEN_EXPIRY" envDefault:"15m"`
		RefreshTokenExpiry time.Duration `env:"REFRESH_TOKEN_EXPIRY" envDefault:"24h"`
	}
}
