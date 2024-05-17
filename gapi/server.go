package gapi

import (
	"fmt"
	db "simple-bank/db/sqlc"
	"simple-bank/pb"
	"simple-bank/token"
	"simple-bank/utils"
)

type Server struct {
	pb.UnimplementedSimpleBankServer
	config     utils.Config
	store      db.Store
	tokenMaker token.Maker
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
