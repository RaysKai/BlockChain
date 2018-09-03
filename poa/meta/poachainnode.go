package meta

import (
	"github.com/linkchain/meta/block"
	"github.com/linkchain/common/math"
)

type POAChainNode struct {
	CurentHash math.Hash
	PrevHash math.Hash
	Height uint32
}

func NewPOAChainNode(block POABlock) POAChainNode {
	return POAChainNode{
		CurentHash: block.GetBlockID().(math.Hash),
		PrevHash: block.GetPrevBlockID().(math.Hash),
		Height: block.GetHeight()}
}

func (bn *POAChainNode) GetNodeHeight() uint32 {
	return bn.Height
}

func (bn *POAChainNode) GetNodeHash() block.IBlockID {
	return bn.CurentHash
}

func (bn *POAChainNode) CheckPrev(prevNode POAChainNode) bool  {
	return bn.PrevHash.IsEqual(&prevNode.CurentHash)
}


