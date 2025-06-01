package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	psqlMigrator "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type MigrationConfig struct {
	RunMigration   bool
	MigrationsPath string
}

type SSLConfig struct {
	SSLMode     string
	SSLRootCert string
}

type ConnectionConfig struct {
	PingTimeout  time.Duration
	MaxIdleTime  time.Duration
	MaxLifetime  time.Duration
	MaxOpenConns int
	MaxIdleConns int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	DBName   string
	Username string
	Password string
	ConnectionConfig
	SSLConfig
	MigrationConfig
}

func (c DatabaseConfig) Validate() error {
	if c.Host == "" {
		return errors.New("database host is required")
	}
	if c.Port == 0 {
		return errors.New("database port is required")
	}
	if c.DBName == "" {
		return errors.New("database name is required")
	}
	if c.Username == "" {
		return errors.New("database username is required")
	}
	if c.Password == "" {
		return errors.New("database password is required")
	}
	if c.PingTimeout == time.Duration(0) {
		return errors.New("database ping timeout is missing")
	}
	if c.MaxIdleTime == time.Duration(0) {
		return errors.New("database max idle time is missing")
	}
	if c.MaxLifetime == time.Duration(0) {
		return errors.New("database max lifetime is missing")
	}
	if c.SSLMode == "" {
		return errors.New("ssl mode is required")
	}
	if c.RunMigration && c.MigrationsPath == "" {
		return errors.New("migration config: migrations path is required")
	}
	if c.SSLMode != "disabled" && c.SSLRootCert == "" {
		return errors.New("ssl root cert is required")
	}
	return nil
}

func InitDatabase(ctx context.Context, conf DatabaseConfig) (CommonRepository, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}

	dbConn := fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		conf.Host, conf.Port, conf.DBName, conf.Username, conf.Password, conf.SSLMode)

	if conf.SSLMode != "disabled" {
		dbConn += fmt.Sprintf(" sslrootcert=%s", conf.SSLRootCert)
	}

	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetConnMaxIdleTime(conf.MaxIdleTime)
	db.SetConnMaxLifetime(conf.MaxLifetime)

	ctx, cancel := context.WithTimeout(ctx, conf.PingTimeout)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	c := &commonRepository{
		client: db,
		dbName: conf.DBName,
	}

	if conf.MigrationConfig.RunMigration {
		err = c.migrateDatabase(conf.MigrationConfig.MigrationsPath)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *commonRepository) migrateDatabase(path string) error {
	driver, err := psqlMigrator.WithInstance(c.client, &psqlMigrator.Config{})
	if err != nil {
		return fmt.Errorf("migration driver creation failed: %s", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+path, c.dbName, driver)
	if err != nil {
		return fmt.Errorf("migration instance creation failed: %s", err)
	}

	err = m.Migrate(findLastMigrationVersion(path))
	switch {
	case err == migrate.ErrNoChange:
		return nil
	case err != nil:
		return fmt.Errorf("migration failed: %s", err)
	default:
		return nil
	}
}

func findLastMigrationVersion(path string) uint {
	files, err := os.ReadDir(path)
	if err != nil {
		return 0
	}

	var migrations []int

	for _, f := range files {
		parts := strings.Split(f.Name(), "_")
		version, err := strconv.Atoi(parts[0])
		if err == nil {
			migrations = append(migrations, version)
		}
	}

	if len(migrations) == 0 {
		return 0
	}

	sort.Ints(migrations)

	return uint(migrations[len(migrations)-1])
}
