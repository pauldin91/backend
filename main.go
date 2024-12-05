package main

import (
	"context"
	"log"
	"net"

	"github.com/pauldin91/backend/api"
	db "github.com/pauldin91/backend/db/sqlc"
	"github.com/pauldin91/backend/gapi"
	"github.com/pauldin91/backend/pb"
	"github.com/pauldin91/backend/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

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
	runGrpcServer(cfg, store)
	//runGinServer(cfg, store)

}

func runGrpcServer(cfg utils.Config, store db.Store) {
	server, err := gapi.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		log.Fatal("Could not create listener")
	}

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot start gRPC server")
	}

}

func runGinServer(cfg utils.Config, store db.Store) {
	server, err := api.NewServer(cfg, store)
	if err != nil {
		log.Fatal(err)
	}
	err = server.Start(cfg.HTTPServerAddress)

	if err != nil {
		log.Fatal("Could not start server")
	}
}
