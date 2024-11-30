package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testQueries *Queries

const (
	dbDriver = "postgres"
	dbSource = "postgres://backend:backend@localhost:5433/backend?sslmode=disable"
)

var testStore Store

func TestMain(m *testing.M) {
	//config, err := util.LoadConfig("../..")
	//if err != nil {
	//	log.Fatal("cannot load config:", err)
	//}

	connPool, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testStore = NewStore(connPool)
	os.Exit(m.Run())
}
