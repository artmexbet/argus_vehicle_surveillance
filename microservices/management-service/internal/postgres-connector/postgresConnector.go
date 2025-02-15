package postgres_connector

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	User     string `yaml:"user" env:"USER" env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"postgres"`
	DBName   string `yaml:"dbname" env:"DBNAME" env-default:"postgres"`
}

type PostgresConnector struct {
	cfg  *Config
	conn *pgxpool.Pool
}

func New(cfg *Config) (*PostgresConnector, error) {
	ctx := context.Background()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	conf, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}
	conf.MaxConns = 80
	conn, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		return nil, err
	}
	return &PostgresConnector{cfg: cfg, conn: conn}, nil
}

func (p *PostgresConnector) Connection() *pgxpool.Pool {
	return p.conn
}

func (p *PostgresConnector) Close() {
	p.conn.Close()
}
