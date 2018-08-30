package meta

import (
	"github.com/linkchain/common/math"
)

const(
	StateLess = 0x00000000 //0
	OnMainChain = 0x00000001 //1
	ChainIDCheck = 0xffffffff
)

type BlockIndex struct {
	Block POABlock	//the data of block
	Status uint32		//the state of block 0x00000000(default):need to deal ;0x000000001 is on main chain ;0x000000010 is on side(1) chain
	Blockhash math.Hash //the hash of block
}

func NewBlockIndex(block POABlock,status uint32) BlockIndex  {
	return BlockIndex{Block:block,Status:status,Blockhash:(block.GetBlockID().(math.Hash))}
}

func (blockIndex *BlockIndex) isOnMainChain() bool  {
	return blockIndex.Status == 1
}

func (blockIndex *BlockIndex) isOnChain() bool  {
	return blockIndex.isOnMainChain() || blockIndex.isOnSideChain()
}

func (blockIndex *BlockIndex) isOnSideChain() bool  {
	return blockIndex.Status > 1
}
