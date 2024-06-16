//go:build wireinject
// +build wireinject

package di

import (
	http "awesomeProject/pkg/api"
	"awesomeProject/pkg/api/handler"
	"awesomeProject/pkg/config"
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/service"
	"github.com/google/wire"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	wire.Build(db.ConnectToBD, service.NewService, handler.NewHandler, http.NewServerHTTP)

	return &http.ServerHTTP{}, nil
}
