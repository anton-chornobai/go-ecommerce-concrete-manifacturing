package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/anton-chornobai/beton.git/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

func OpenPostgre(connStr string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connStr)

	if err != nil {
		return nil, fmt.Errorf("failed to open db, %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return db, nil
}

func GetDBConnStr(cfg config.DBConfig) string {
	return fmt.Sprintf(
		"postgresql://%s:%d@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLmode,
	)
}
