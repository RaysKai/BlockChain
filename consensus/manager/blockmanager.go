package manager

import (
	"github.com/linkchain/meta/block"
	"github.com/linkchain/common"
)

type BlockManager interface {
	common.IService
	BlockBaseManager
	BlockPoolManager
	BlockValidator

	ProcessBlock(block block.IBlock)
}

type BlockBaseManager interface {
	NewBlock() block.IBlock
	GetGensisBlock() block.IBlock
}
/**
	BlockPoolManager
	manager block pool.the block at pool is side chain block or no-parent block
 */
type BlockPoolManager interface{
	GetBlockByID(hash block.IBlockID) (block.IBlock,error)
	GetBlockByHeight(height uint32) ([]block.IBlock,error)
	AddBlock(block block.IBlock) error
	AddBlocks(block []block.IBlock) error
	RemoveBlock(hash block.IBlockID) error
}

type BlockValidator  interface {
	CheckBlock(block block.IBlock) bool
}