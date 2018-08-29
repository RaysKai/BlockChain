package poa

import (
	"time"

	"github.com/linkchain/common/util/log"
	"github.com/linkchain/poa/poamanager"
	"github.com/linkchain/meta/block"
	"github.com/linkchain/common/math"
	"github.com/linkchain/meta/tx"

	poameta "github.com/linkchain/poa/meta"
)

type POAValidator struct {

}

func (validator *POAValidator) CheckBlock(block block.IBlock) bool {
	log.Info("poa CheckBlock ...")
	poamanager.GetInstance().ChainMgr.AddBlock(block)
	log.Info("poa Add a Block","block",block)
	log.Info("poa GetBestHeader","blockhash", poamanager.GetInstance().ChainMgr.GetBestBlock().GetBlockID().GetString())
	txs := []poameta.POATransaction{}
	block1 := &poameta.POABlock{
		Header:poameta.POABlockHeader{Version:1, PrevBlock:math.Hash{},MerkleRoot:math.Hash{},Timestamp:time.Unix(1401292357, 0),Difficulty:0x207fffff,Nonce:0,Extra:nil},
		TXs:txs,
	}
	log.Info("poa Add a Block","block",block1)
	poamanager.GetInstance().ChainMgr.AddBlock(block1)
	log.Info("poa GetBestHeader","blockhash",poamanager.GetInstance().ChainMgr.GetBestBlock().GetBlockID().GetString())

	return true
}

func (validator *POAValidator) CheckTx(tx tx.ITransaction) bool {
	log.Info("poa CheckTx ...")
	return true
}