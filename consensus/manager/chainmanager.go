package manager

import (
	"github.com/linkchain/meta"
	"github.com/linkchain/common"
)


type ChainReader interface{

	GetBestHeader() *meta.BlockHeader

	GetMainChain() *meta.BlockHeader

	GetBlockByHash(h *meta.Hash) meta.Block

	GetBlockByHeight(height uint32) meta.Block

	GetBlockChainInfo() string
}

type ChainWriter interface{

	AddBlock(block meta.Block)

	UpdateChain() bool
}

type ChainManager interface{
	common.IService
	ChainWriter
	ChainReader
}
