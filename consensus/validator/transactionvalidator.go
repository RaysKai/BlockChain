package validator

import "github.com/linkchain/meta"

type TransactionVlidator interface {
	CheckTx(tx meta.Transaction) bool
}
