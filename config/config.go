package config

import (
	"context"
	"time"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var Cfg Config

func LoadConfig(ctx context.Context) {
	if err := godotenv.Load(); err != nil {
		log.Error(err.Error())
	}

	err := env.Parse(&Cfg)
	if err != nil {
		log.Error(err.Error())
	}
}

type Config struct {
	Environment string `env:"ENVIRONMENT,required"`
	ServiceName string `env:"SERVICE_NAME" envDefault:""`
	Timezone    string `env:"TIMEZONE,required" envDefault:"Asia/Jakarta"`

	HttpPort int    `env:"API_PORT" envDefault:"8080"`
	Host     string `env:"GRPC_HOST" envDefault:"localhost"`

	Database struct {
		Host        string        `env:"DB_HOST" envDefault:"localhost"`
		Port        int           `env:"DB_PORT" envDefault:"5433"`
		DBName      string        `env:"DB_NAME" envDefault:"shortina"`
		Username    string        `env:"DB_USERNAME"`
		Password    string        `env:"DB_PASSWORD"`
		PingTimeout time.Duration `env:"DB_TIMEOUT" envDefault:"5s"`
		Migration   struct {
			RunMigration bool   `env:"DB_RUN_MIGRATION" envDefault:"true"`
			Path         string `env:"DB_MIGRATION_PATH" envDefault:"storage/postgres/migrations"`
		}
	}

	Redis struct {
		Host string `env:"REDIS_HOST" envDefault:"localhost"`
		Port int    `env:"REDIS_PORT" envDefault:"6379"`
	}

	Token struct {
		SecretKey          string        `env:"TOKEN_SECRET_KEY" envDefault:"gomugomuno"`
		AccessTokenExpiry  time.Duration `env:"ACCESS_TOKEN_EXPIRY" envDefault:"15m"`
		RefreshTokenExpiry time.Duration `env:"REFRESH_TOKEN_EXPIRY" envDefault:"24h"`
	}
}
