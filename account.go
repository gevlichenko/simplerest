package main

import (
	"errors"
	"sync"
)

type Account struct {
	mutex     sync.Mutex
	Id        uint
	FirstName string
	LastName  string
	balance   float64
}

func NewAccount(id uint, firstName string, lastName string, balance float64) (*Account, error) {
	if firstName == "" || lastName == "" {
		return nil, errors.New("first or last name isn't present")
	}
	account := new(Account)
	account.Id = id
	account.FirstName = firstName
	account.LastName = lastName
	account.balance = balance
	return account, nil
}
func (account *Account) GetBalance() float64 {
	return account.balance
}
func (account *Account) SetBalance(balance float64) {
	account.balance = balance
}
func (account *Account) Lock() {
	account.mutex.Lock()
}
func (account *Account) Unlock() {
	account.mutex.Unlock()
}
