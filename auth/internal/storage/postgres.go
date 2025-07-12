package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func MustConnect(url string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalf("DB connection error: %s", err)
	}
	return pool
}
