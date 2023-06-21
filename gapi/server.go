package gapi

import (
	"fmt"

	db "github.com/thetherington/simplebank/db/sqlc"
	"github.com/thetherington/simplebank/pb"
	"github.com/thetherington/simplebank/token"
	"github.com/thetherington/simplebank/util"
	"github.com/thetherington/simplebank/worker"
)

// Server serves HTTP requests for our banking service
type Server struct {
	config          util.Config
	store           db.Store
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor

	pb.UnimplementedSimpleBankServer
}

func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:           store,
		tokenMaker:      tokenMaker,
		config:          config,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
