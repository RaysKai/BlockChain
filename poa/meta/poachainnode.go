package meta

import (
	"github.com/linkchain/meta/block"
	"github.com/linkchain/common/math"
	"strconv"
)

type POAChainNode struct {
	curentHash math.Hash
	prevHash math.Hash
	height uint32
}

func NewPOAChainNode(block block.IBlock) POAChainNode {
	return POAChainNode{
		curentHash: block.GetBlockID().(math.Hash),
		prevHash: block.GetPrevBlockID().(math.Hash),
		height: block.GetHeight()}
}

func (bn *POAChainNode) GetNodeHeight() uint32 {
	return bn.height
}

func (bn *POAChainNode) GetNodeHash() block.IBlockID {
	return bn.curentHash
}

func (bn *POAChainNode) GetPrevHash() block.IBlockID {
	return bn.prevHash
}

func (bn *POAChainNode) CheckPrev(prevNode POAChainNode) bool  {
	return bn.prevHash.IsEqual(&prevNode.curentHash)
}

func (bn *POAChainNode) IsEuqal(checkNode POAChainNode) bool {
	return bn.curentHash.IsEqual(&checkNode.curentHash)
}

func (bn *POAChainNode) GetString() string {
	str := string("height:") + strconv.Itoa(int(bn.height)) + " "
	str += string("currentHash:") + bn.curentHash.GetString() + " "
	str += string("prev:") + bn.prevHash.GetString()

	return str
}

