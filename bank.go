package main

import (
	"errors"
	"sync"
)

type Bank struct {
	mutex sync.RWMutex
	accounts map[uint]*Account
	lastId   uint
}

func NewBank() *Bank {
	bank := new(Bank)
	bank.accounts = make(map[uint]*Account)
	return bank
}

func (bank *Bank) CreateAccount(firstName string, lastName string, balance float64) (uint, error) {
	id := bank.getNewId()
	account, err := NewAccount(id, firstName, lastName, balance)
	if err != nil {
		return 0, err
	}
	bank.mutex.Lock()
	bank.accounts[id] = account
	bank.mutex.Unlock()
	return id, nil
}
func (bank *Bank) GetBalanceById(id uint) (float64, error) {
	bank.mutex.RLock()
	if val, ok := bank.accounts[id]; ok {
		return val.GetBalance(), nil
	}
	bank.mutex.RUnlock()
	return 0, errors.New("account id isn't valid")
}
func (bank *Bank) MoveCash(senderId, recipientId uint, amount float64) error {
	var sender, recipient *Account
	var ok bool
	bank.mutex.RLock()
	if sender, ok = bank.accounts[senderId]; !ok {
		return errors.New("sender id isn't valid")
	}
	if recipient, ok = bank.accounts[recipientId]; !ok {
		return errors.New("recipient id isn't valid")
	}
	sender.Lock()
	recipient.Lock()
	if amount > sender.GetBalance() {
		return errors.New("balance too low")
	}
	sender.SetBalance(sender.GetBalance() - amount)
	recipient.SetBalance(recipient.GetBalance() + amount)
	recipient.Unlock()
	sender.Unlock()
	bank.mutex.RUnlock()
	return nil
}
func (bank *Bank) getNewId() uint {
	bank.mutex.Lock()
	bank.lastId = bank.lastId + 1
	bank.mutex.Unlock()
	return bank.lastId
}
