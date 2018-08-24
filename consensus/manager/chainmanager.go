package manager

import "github.com/linkchain/meta"

type ChainManager interface{}

type ChainReader interface{

	GetBestHeader() *meta.BlockHeader

	GetMainChain() *meta.BlockHeader

	GetBlockByHash(h *meta.Hash) meta.Block

	GetBlockByHeight(height uint32) meta.Block

	GetBlockChainInfo() string
}
