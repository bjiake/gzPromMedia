package interfaces

import (
	"awesomeProject/pkg/domain/account"
	"context"
)

type AccountRepository interface {
	Migrate(ctx context.Context) error
	Registration(ctx context.Context, newAccount account.Registration) (*account.Info, error)
	Login(ctx context.Context, acc account.Login) (int64, error)
	Put(ctx context.Context, id int64, updateAcc *account.Account) (*account.Info, error)
	Get(ctx context.Context, id int64) (*account.Account, error)
	GetAll(ctx context.Context) ([]account.Account, error)
	Delete(ctx context.Context, id int64) error
}
