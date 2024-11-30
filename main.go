package main

import (
	"backend/api"
	db "backend/db/sqlc"
	"backend/utils"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := pgxpool.New(context.Background(), cfg.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(cfg.HTTPServerAddress)

	if err != nil {
		log.Fatal("Could not start server")
	}

}
