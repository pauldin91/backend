package db

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v5"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgres://backend:backend@localhost:5433/backend?sslmode=disable"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dbSource)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)
	testQueries = New(conn)
	m.Run()
}
