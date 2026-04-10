package database

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"time"

	"timetable-homework-tgbot/internal/infrastracture/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

const (
	postgresDialect     = "postgres"
	migrationsDirectory = "migrations"
	minConnections      = 2
	maxConnections      = 10
	maxConnIdleTime     = 100 * time.Millisecond
	maxConnLifetime     = time.Second
)

var (
	ErrCreatingDatabaseFromURL = errors.New("migrating database error")
	ErrMigrating               = errors.New("migrating database error")
	ErrCreatingConnectionPool  = errors.New("creating connection db error")
)

type DB struct {
	db *pgxpool.Pool
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func NewDB(config config.PostgresConfig) (*DB, error) {
	dsn := createDataBaseURLFromConfig(config)

	connConf, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Join(err, ErrCreatingDatabaseFromURL)
	}

	if err = migrate(connConf); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrMigrating, err)
	}

	poolConf, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Join(err, ErrCreatingDatabaseFromURL)
	}
	poolConf.MinConns = minConnections
	poolConf.MaxConns = maxConnections
	poolConf.MaxConnLifetime = maxConnLifetime
	poolConf.MaxConnIdleTime = maxConnIdleTime

	ctx := context.Background()
	pool, err := pgxpool.NewWithConfig(ctx, poolConf)
	if err != nil {
		return nil, errors.Join(err, ErrCreatingConnectionPool)
	}

	return &DB{db: pool}, nil
}

func migrate(cfg *pgx.ConnConfig) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(postgresDialect); err != nil {
		return errors.Join(err, ErrMigrating)
	}
	db := stdlib.OpenDB(*cfg)
	defer func() { _ = db.Close() }()

	if err := goose.Up(db, migrationsDirectory); err != nil {
		return errors.Join(err, ErrMigrating)
	}
	return nil
}

func createDataBaseURLFromConfig(cfg config.PostgresConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
}

func (db *DB) Close() {
	if db.db != nil {
		db.db.Close()
	}
}
