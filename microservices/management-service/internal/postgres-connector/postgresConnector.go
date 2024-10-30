package postgres_connector

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
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
	conn *pgx.Conn
}

func New(cfg *Config) (*PostgresConnector, error) {
	ctx := context.Background()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresConnector{cfg: cfg, conn: conn}, nil
}

func (p *PostgresConnector) Connection() *pgx.Conn {
	return p.conn
}

func (p *PostgresConnector) Close() {
	p.conn.Close(context.Background())
}
