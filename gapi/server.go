package gapi

import (
	"fmt"

	db "github.com/pauldin91/backend/db/sqlc"
	pb "github.com/pauldin91/backend/pb"
	"github.com/pauldin91/backend/token"
	"github.com/pauldin91/backend/utils"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(cfg utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker %w", err)
	}

	server := &Server{
		config:     cfg,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
