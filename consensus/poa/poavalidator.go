package poa

import (
	"time"

	"github.com/linkchain/meta"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/consensus/poa/poamanager"
)

type POAValidator struct {

}

func (validator *POAValidator) CheckBlock(block meta.Block) bool {
	log.Info("poa CheckBlock ...")
	log.Info("poa GetBestHeader","blockhash", poamanager.GetInstance().ChainMgr.GetBestHeader().GetHash().String())
	txs := []meta.Transaction{}
	block1 := meta.Block{
		Header:meta.BlockHeader{Version:1, PrevBlock:meta.Hash{},MerkleRoot:meta.Hash{},Timestamp:time.Unix(1401292357, 0),Difficulty:0x207fffff,Nonce:0,Extra:nil},
		TXs:txs,
	}
	log.Info("poa Add a Block","block",block1)
	poamanager.GetInstance().ChainMgr.AddBlock(block1)
	log.Info("poa GetBestHeader","blockhash",poamanager.GetInstance().ChainMgr.GetBestHeader().GetHash().String())
	return true
}

func (validator *POAValidator) CheckTx(tx meta.Transaction) bool {
	log.Info("poa CheckTx ...")
	return true
}