package manager

import (
	"github.com/linkchain/meta/block"
	"github.com/linkchain/common"
)

type BlockManager interface {
	common.IService
	//todo BlockPoolManager
	BlockValidator
	//todo NewBlock() block.IBlock

	ProcessBlock(block block.IBlock)
}

/**
	BlockPoolManager
	manager block pool.the block at pool is side chain block or no-parent block
 */
type BlockPoolManager interface{
	GetBlockByID(hash block.IBlockID) (block.IBlock,error)
	GetBlockByHeight(height uint32) (block.IBlock,error)
	FindLongestChain() (uint32,uint32,error)/**return the start height of the longest side chain,the end height of the longest side chain,error**/
	AddBlock(block block.IBlock) error
	AddBlocks(block []block.IBlock) error
	RemoveBlock(hash block.IBlockID) error
}

type BlockValidator  interface {
	CheckBlock(block block.IBlock) bool
}