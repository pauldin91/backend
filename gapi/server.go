package gapi

import (
	db "github.com/pauldin91/backend/db/sqlc"
	"github.com/pauldin91/backend/token"
	"github.com/pauldin91/backend/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
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
	if _, ok := binding.Validator.Engine().(*validator.Validate); ok {
		//		v.RegisterValidation("currency", validCurrency)
	}

	return server, nil
}
