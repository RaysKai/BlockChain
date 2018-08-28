package manager

import (
	"hash"
	"github.com/linkchain/meta"
)

type TransactionManager interface{
	GetTransaction(txid hash.Hash) meta.Transaction
}

