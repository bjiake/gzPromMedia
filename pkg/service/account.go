package service

import (
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/domain/account"
	"context"
	"log"
	"regexp"
	"strings"
)

func (s *service) Subscribe(ctx context.Context, id string, idSub string) error {
	idInt, err := s.checkIdParam(id)
	if err != nil {
		return err
	}
	idSubInt, err := s.checkIdParam(idSub)
	if err != nil {
		return err
	}
	currentAccount, err := s.rAccount.Get(ctx, idInt)
	if err != nil {
		return err
	}
	subAccount, err := s.rAccount.Get(ctx, idSubInt)
	if err != nil {
		return err
	}
	subAccount.Subscribe(currentAccount)

	_, err = s.rAccount.Put(ctx, idInt, subAccount)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) UnSubscribe(ctx context.Context, id string, idSub string) error {
	idInt, err := s.checkIdParam(id)
	if err != nil {
		return err
	}
	idSubInt, err := s.checkIdParam(idSub)
	if err != nil {
		return err
	}
	currentAccount, err := s.rAccount.Get(ctx, idInt)
	if err != nil {
		return err
	}
	subAccount, err := s.rAccount.Get(ctx, idSubInt)
	if err != nil {
		return err
	}
	subAccount.Unsubscribe(currentAccount)

	_, err = s.rAccount.Put(ctx, idInt, subAccount)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) Registration(ctx context.Context, userId string, newAccount account.Registration) (*account.Info, error) {
	if userId != "" {
		return nil, db.ErrAuthorize
	}
	//check valid data
	if !isValidAccountRegister(newAccount) {
		log.Println("invalid data")
		return nil, db.ErrValidate
	}

	result, err := s.rAccount.Registration(ctx, newAccount)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) Login(ctx context.Context, acc account.Login) (int64, error) {
	id, err := s.rAccount.Login(ctx, acc)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *service) PutAccount(ctx context.Context, id string, updateAcc *account.Account) (*account.Info, error) {
	idInt, err := s.checkIdParam(id)
	if err != nil {
		return nil, err
	}

	tempAccount := &account.Registration{
		FirstName: updateAcc.FirstName,
		LastName:  updateAcc.LastName,
		Email:     updateAcc.Email,
		Password:  updateAcc.Password,
	}

	// Валидация обновляемых данных
	if !isValidAccountRegister(*tempAccount) {
		return nil, db.ErrValidate
	}

	result, err := s.rAccount.Put(ctx, idInt, updateAcc)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *service) GetAccount(ctx context.Context, id string) (*account.Account, error) {
	idInt, err := s.checkIdParam(id)
	if err != nil {
		return nil, err
	}
	result, err := s.rAccount.Get(ctx, idInt)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *service) DeleteAccount(ctx context.Context, id string) error {
	idInt, err := s.checkIdParam(id)
	if err != nil {
		return err
	}

	err = s.rAccount.Delete(ctx, idInt)
	if err != nil {
		return err
	}
	return nil
}

func isValidAccountRegister(newAccountRegistration account.Registration) bool {
	if newAccountRegistration.FirstName == "" || strings.TrimSpace(newAccountRegistration.FirstName) == "" {
		return false
	}

	if newAccountRegistration.LastName == "" || strings.TrimSpace(newAccountRegistration.LastName) == "" {
		return false
	}

	if newAccountRegistration.Email == "" || strings.TrimSpace(newAccountRegistration.Email) == "" || !isValidEmail(newAccountRegistration.Email) {
		return false
	}

	if newAccountRegistration.Password == "" || strings.TrimSpace(newAccountRegistration.Password) == "" {
		return false
	}

	return true
}

func isValidEmail(email string) bool {
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, email)
	return match
}
