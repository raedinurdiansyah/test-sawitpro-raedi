package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/SawitProRecruitment/UserService/config"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()
	e.Logger.SetLevel(log.INFO)
	cfg := config.GetConfig()
	var server generated.ServerInterface = newServer(cfg)

	usersGroup := e.Group("/users")
	usersGroup.Use(handler.JWTMiddleware(*cfg))
	usersGroup.GET("", server.GetUserProfile)
	usersGroup.PUT("", server.UpdateUser)

	generated.RegisterHandlers(e, server)

	// Start server
	go func(port uint16) {
		if err := e.Start(fmt.Sprintf(":%v", cfg.Port)); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}(cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func newServer(config *config.Config) *handler.Server {
	var repo repository.RepositoryInterface = repository.NewRepository(
		repository.NewRepositoryOptions{
			Dsn: config.DatabaseURL,
		},
	)
	opts := handler.NewServerOptions{
		Repository: repo,
		Config:     *config,
	}
	return handler.NewServer(opts)
}
