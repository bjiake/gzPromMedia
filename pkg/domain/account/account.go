package account

import (
	"sync"
	"time"
)

type Account struct {
	ID             int64     `json:"id" validate:"required" `
	FirstName      string    `json:"firstName" validate:"required" `
	LastName       string    `json:"lastName" validate:"required" `
	BirthDate      time.Time `json:"birthDate" validate:"required" `
	Email          string    `json:"email" validate:"required" `
	Password       string    `json:"password" validate:"required" `
	SubscribersIds []int64   `json:"subscribers"`
	mut            sync.RWMutex
}

func (a *Account) Subscribe(sub *Account) {
	a.mut.Lock()
	defer a.mut.Unlock()

	a.SubscribersIds = append(a.SubscribersIds, sub.ID)
}

func (a *Account) Unsubscribe(sub *Account) {
	a.mut.Lock()
	defer a.mut.Unlock()
	for i, s := range a.SubscribersIds {
		if s == sub.ID {
			a.SubscribersIds = append(a.SubscribersIds[:i], a.SubscribersIds[i+1:]...)
		}
	}
}

// IsBirthday Проверяет день рождение на сегодня
func (a *Account) IsBirthday() bool {
	if a.BirthDate.Day() == time.Now().Day() && a.BirthDate.Month() == time.Now().Month() {
		return true
	}

	return false
}

// IsBornOnLeapYear На случай, если дата было 29 февраля
func (a *Account) IsBornOnLeapYear() bool {
	return a.BirthDate.Day() == 29 && a.BirthDate.Month() == time.February
}

type Info struct {
	ID          int64     `json:"id" validate:"required" `
	FirstName   string    `json:"firstName" validate:"required" `
	LastName    string    `json:"lastName" validate:"required" `
	BirthDate   time.Time `json:"birthDate" validate:"required" `
	Email       string    `json:"email" validate:"required" `
	Subscribers []int64   `json:"subscribers"`
}

type Login struct {
	Email    string `json:"email" validate:"required" `
	Password string `json:"password" validate:"required" `
}

type Registration struct {
	FirstName string    `json:"firstName" validate:"required" `
	LastName  string    `json:"lastName" validate:"required" `
	BirthDate time.Time `json:"birthDate" validate:"required" `
	Email     string    `json:"email" validate:"required" `
	Password  string    `json:"password" validate:"required" `
}
