package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/viper"
)

const (
	defaultMigrationsPath = "internal/db/migrations"
)

type Config struct {
	DatabaseURL    string
	MigrationsPath string
}

type Database struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (d Database) URL() string {
	sslMode := d.SSLMode
	if sslMode == "" {
		sslMode = "disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(d.User),
		url.QueryEscape(d.Password),
		d.Host,
		d.Port,
		d.Name,
		sslMode,
	)
}

func Load(dir string) (Config, error) {
	v := viper.New()
	if dir != "" {
		v.AddConfigPath(dir)
	}
	v.SetConfigType("env")
	v.SetConfigName(".env")

	keys := []string{
		"MIGRATION_URL",
		"MIGRATIONS_PATH",
		"DB_HOST",
		"DB_PORT",
		"DB_USERNAME",
		"DB_PASSWORD",
		"DB_DATABASE",
		"DB_SSLMODE",
	}
	for _, key := range keys {
		if err := v.BindEnv(key); err != nil {
			return Config{}, err
		}
	}

	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return Config{}, err
		}
	}

	cfg := Config{
		DatabaseURL:    firstNonEmpty(v.GetString("MIGRATION_URL"), os.Getenv("MIGRATION_URL")),
		MigrationsPath: firstNonEmpty(v.GetString("MIGRATIONS_PATH"), os.Getenv("MIGRATIONS_PATH"), defaultMigrationsPath),
	}

	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = Database{
			Host:     firstNonEmpty(v.GetString("DB_HOST"), os.Getenv("DB_HOST")),
			Port:     firstNonEmpty(v.GetString("DB_PORT"), os.Getenv("DB_PORT"), "5432"),
			User:     firstNonEmpty(v.GetString("DB_USERNAME"), os.Getenv("DB_USERNAME")),
			Password: firstNonEmpty(v.GetString("DB_PASSWORD"), os.Getenv("DB_PASSWORD")),
			Name:     firstNonEmpty(v.GetString("DB_DATABASE"), os.Getenv("DB_DATABASE")),
			SSLMode:  firstNonEmpty(v.GetString("DB_SSLMODE"), os.Getenv("DB_SSLMODE")),
		}.URL()
	}

	return cfg, nil
}

func New(databaseURL, migrationsPath string) (Config, error) {
	cfg := Config{
		DatabaseURL:    databaseURL,
		MigrationsPath: migrationsPath,
	}
	return cfg, cfg.Validate()
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.DatabaseURL) == "" {
		return errors.New("database url is required")
	}
	if strings.TrimSpace(c.MigrationsPath) == "" {
		return errors.New("migrations path is required")
	}
	return nil
}

func (c Config) MigrationsSource() string {
	path := strings.TrimSpace(c.MigrationsPath)
	if strings.HasPrefix(path, "file://") {
		return path
	}
	return "file://" + strings.TrimPrefix(path, "/")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
