package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pauldin91/backend/api"
	db "github.com/pauldin91/backend/db/sqlc"
	_ "github.com/pauldin91/backend/doc/statik"
	"github.com/pauldin91/backend/gapi"
	pb "github.com/pauldin91/backend/pb"
	"github.com/pauldin91/backend/utils"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
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
	runDbMigration(cfg.MigrationURL, cfg.DBSource)

	store := db.NewStore(conn)
	go runGatewayServer(cfg, store)
	runGrpcServer(cfg, store)
	//runGinServer(cfg, store)

}

func runDbMigration(migrationUrl string, dbSource string) {
	migration, err := migrate.New(migrationUrl, dbSource)
	if err != nil {
		log.Fatal("failed to create migrations ", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to apply migrations ", err)
	}
	log.Println("Successfully applied migrations")
}

func runGrpcServer(cfg utils.Config, store db.Store) {
	server, err := gapi.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)

	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		log.Fatal("Could not create listener")
	}

	log.Printf("start gPRC server at %s", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot start gRPC server")
	}

}

func runGatewayServer(cfg utils.Config, store db.Store) {
	server, err := gapi.NewServer(cfg, store)
	if err != nil {
		log.Fatal("cannot create server: ", err)
	}
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik ", err)
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", cfg.HTTPServerAddress)
	if err != nil {
		log.Fatal("Could not create listener")
	}

	log.Printf("start gPRC server at %s", listener.Addr().String())

	err = http.Serve(listener, mux)

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
