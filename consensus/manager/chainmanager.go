package manager

import (
	"github.com/linkchain/common"
	"github.com/linkchain/meta/block"
	"github.com/linkchain/common/math"
)


type ChainReader interface{

	GetBestBlock() block.IBlock


	GetBlockByHash(h math.Hash) block.IBlock

	GetBlockByHeight(height uint32) block.IBlock

	GetBlockChainInfo() string

	GetBlockAncestor(block block.IBlock,height uint32) block.IBlock
}

type ChainWriter interface{

	AddBlock(block block.IBlock)

	UpdateChain() bool
}

type ChainManager interface{
	common.IService
	ChainWriter
	ChainReader
}
