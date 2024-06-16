package service

import (
	"awesomeProject/pkg/db"
	accountI "awesomeProject/pkg/repo/account/interface"

	interfaces "awesomeProject/pkg/service/interface"
	"context"
	"strconv"
)

type service struct {
	rAccount accountI.AccountRepository
}

func NewService(accountRepository accountI.AccountRepository) interfaces.ServiceUseCase {
	return &service{
		rAccount: accountRepository,
	}
}

func (s *service) Migrate(ctx context.Context) error {
	if err := s.rAccount.Migrate(ctx); err != nil {
		return err
	}

	return nil
}

func (s *service) checkIdParam(id string) (int64, error) {
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil || idInt <= 0 {
		return 0, db.ErrParamNotFound
	}
	return idInt, nil
}
