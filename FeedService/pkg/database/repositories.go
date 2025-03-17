package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToDB(connStr string) (*pgxpool.Pool, error) {
	db, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}
	return db, nil
} 