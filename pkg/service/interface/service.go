package interfaces

import (
	"awesomeProject/pkg/domain/account"
	"context"
)

type ServiceUseCase interface {
	Migrate(ctx context.Context) error

	//Account
	Registration(ctx context.Context, userId string, newAccount account.Registration) (*account.Info, error)
	Login(ctx context.Context, acc account.Login) (int64, error)
	PutAccount(ctx context.Context, id string, updateAcc *account.Account) (*account.Info, error)
	GetAccount(ctx context.Context, id string) (*account.Account, error)
	DeleteAccount(ctx context.Context, id string) error
	Subscribe(ctx context.Context, id string, idSub string) error
	UnSubscribe(ctx context.Context, id string, idSub string) error
}
