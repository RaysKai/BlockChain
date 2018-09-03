package chain

import (
	"github.com/linkchain/meta/block"
)

type IChain interface {
	//maintain blockchain
	AddNewBlock(block.IBlock)

	Rollback(block.IBlock)

	RollbackAtHeight(int)

	//get blockchain info
	GetLastBlock() block.IBlock

	GetHeight()int

	GetBlockByID(block.IBlockID) block.IBlock

	GetBlockByHeight(int) block.IBlock
}

type IChainGraph interface{

}
