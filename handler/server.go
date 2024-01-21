package handler

import (
	"github.com/SawitProRecruitment/UserService/config"
	"github.com/SawitProRecruitment/UserService/repository"
)

type Server struct {
	Repository repository.RepositoryInterface
	Config     config.Config
}

type NewServerOptions struct {
	Repository repository.RepositoryInterface
	Config     config.Config
}

func NewServer(opts NewServerOptions) *Server {
	return &Server{
		Repository: opts.Repository,
		Config:     opts.Config,
	}
}
