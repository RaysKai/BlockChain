package storage

import (
	"github.com/linkchain/meta/block"
	"github.com/linkchain/meta/tx"
	"github.com/linkchain/meta/chain"
)

type IStroage interface{
	//block
	storeBlock(block block.IBlock)
	loadBlockById(id block.IBlockID) block.IBlock
	loadBlockByHeight(height int) block.IBlock


	//tx
	storeTx(iTx tx.ITx)
	loadTxById(id tx.ITxID) tx.ITx
	loadTxByPeer(peer tx.ITxPeer) []tx.ITx

	//chain info
	storeChain(chain chain.IChain)
	storeChainGraph(graph chain.IChainGraph)
	loadChainGraph()chain.IChainGraph
}