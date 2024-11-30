package main

import (
	"backend/api"
	db "backend/db/sqlc"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	serverAddress = "0.0.0.0:7080"
	dbDriver      = "postgres"
	dbSource      = "postgres://backend:backend@localhost:5433/backend?sslmode=disable"
)

func main() {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("Could not start server")
	}

}
