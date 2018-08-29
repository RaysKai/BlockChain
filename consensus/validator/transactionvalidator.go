package validator

import "github.com/linkchain/meta/tx"

type TransactionVlidator interface {
	CheckTx(tx tx.ITransaction) bool
}
