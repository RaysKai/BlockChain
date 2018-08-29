package manager

import (
	"hash"
	"github.com/linkchain/meta/tx"
)

type TransactionManager interface{
	TransactionValidator
	NewTransaction() tx.ITransaction
	SignTransaction(tx tx.ITransaction) error
	ProcessTx(tx tx.ITransaction)
}

type TransactionPoolManager interface{
	AddTransaction(tx tx.ITransaction) error
	GetTransaction(txid hash.Hash) (tx.ITransaction,error)
	removeTransaction(txid hash.Hash) error

}

type TransactionValidator  interface {
	CheckTx(tx tx.ITransaction) bool
}

