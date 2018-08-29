package manager

import (
	"hash"
	"github.com/linkchain/meta/tx"
)

type TransactionManager interface{
	GetTransaction(txid hash.Hash) tx.ITransaction
}

