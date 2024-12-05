package main

import (
	"github.com/pauldin91/backend/api"
	db "github.com/pauldin91/backend/db/sqlc"
	"github.com/pauldin91/backend/utils"
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
	server, err := api.NewServer(cfg, store)
	if err != nil {
		log.Fatal(err)
	}
	err = server.Start(cfg.HTTPServerAddress)

	if err != nil {
		log.Fatal("Could not start server")
	}

}
