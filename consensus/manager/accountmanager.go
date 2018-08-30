package manager

import "github.com/linkchain/meta/account"

type AccountManager interface{
	//todo AccountPoolManager

	NewAccount() account.IAccount
}

type AccountPoolManager interface{
	AddAccount(iAccount account.IAccount) error
	GetAccount(id account.IAccountID) (account.IAccount,error)
	RemoveAccount(id account.IAccountID) error
	UpdateAccount(iAccount account.IAccount) error
}
