package di

import (
	http "awesomeProject/pkg/api"
	"awesomeProject/pkg/api/handler"
	"awesomeProject/pkg/config"
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/repo/account"
	"awesomeProject/pkg/service"
	"context"
)

func InitializeAPI(cfg config.Config) (*http.ServerHTTP, error) {
	bd, err := db.ConnectToBD(cfg)
	if err != nil {
		return nil, err
	}
	// Repository
	accountRepository := account.NewAccountDataBase(bd)

	//service - logic
	userService := service.NewService(accountRepository)

	// Init Migrate
	err = userService.Migrate(context.Background())
	if err != nil {
		return nil, err
	}

	userHandler := handler.NewHandler(userService)
	serverHTTP := http.NewServerHTTP(userHandler)

	return serverHTTP, nil
}
