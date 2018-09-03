package manager

import (
	"hash"
	"github.com/linkchain/meta/tx"
)

type TransactionManager interface{
	TransactionValidator
	NewTransaction() tx.ITx
	SignTransaction(tx tx.ITx) error
	ProcessTx(tx tx.ITx)
}

type TransactionPoolManager interface{
	AddTransaction(tx tx.ITx) error
	GetTransaction(txid hash.Hash) (tx.ITx,error)
	removeTransaction(txid hash.Hash) error

}

type TransactionValidator  interface {
	CheckTx(tx tx.ITx) bool
}

