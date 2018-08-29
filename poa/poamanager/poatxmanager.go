package poamanager

import (
	"hash"
	"github.com/linkchain/meta/tx"
)

type POATxManager struct {
}

/*func (m *POATxManager) GetTransaction(txid hash.Hash) meta.Transaction  {
	return meta.Transaction{}
}*/

func (m *POATxManager) GetTransaction(txid hash.Hash) tx.ITransaction  {
	return nil
}
