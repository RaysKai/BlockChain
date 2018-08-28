package poamanager

import (
	"time"

	"github.com/linkchain/meta"
)

type POABlockManager struct {

}

func (m *POABlockManager) GetBlockDetail() meta.Block {
	txs := []meta.Transaction{}
	block := meta.Block{
		Header:meta.BlockHeader{Version:1, PrevBlock:meta.Hash{},MerkleRoot:meta.Hash{},Timestamp:time.Unix(1401292357, 0),Difficulty:0x207fffff,Nonce:0,Extra:nil},
		TXs:txs,
	}
	return block
}