package poamanager

import (
	"hash"
	"github.com/linkchain/meta"
)

type POATxManager struct {
}

func (m *POATxManager) GetTransaction(txid hash.Hash) meta.Transaction  {
	return meta.Transaction{}
}
