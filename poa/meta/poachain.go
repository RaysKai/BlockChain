package meta

import (
	"github.com/linkchain/meta/block"
	"errors"
)

type POAChain struct {
	ChainNode []POABlock
}

func NewPOAChain(startNode block.IBlock,endNode block.IBlock) POAChain {
	chainNode := make([]POABlock,0)
	chainNode = append(chainNode,*(startNode.(*POABlock)))
	chainNode = append(chainNode,*(endNode.(*POABlock)))
	return POAChain{ChainNode:chainNode}
}


func (bc *POAChain) AddNewBlock(block block.IBlock) {
	bc.ChainNode = append(bc.ChainNode,*(block.(*POABlock)))
}

/**invalidate block by block*/
func (bc *POAChain) Rollback(block.IBlock) {

}

/**invalidate block by height*/
func (bc *POAChain) RollbackAtHeight(int) {

}

func (bc *POAChain) GetLastBlock() block.IBlock {
	return &bc.ChainNode[len(bc.ChainNode)-1]
}

func (bc *POAChain) GetHeight() uint32  {
	return bc.GetLastBlock().GetHeight()
}

func (bc *POAChain) GetBlockByID(block.IBlockID) block.IBlock {
	//todo need to sorage
	return nil
}

func (bc *POAChain) GetBlockByHeight(int) block.IBlock {
	//todo need to sorage
	return nil
}

func (bc *POAChain) UpdateChainTop(topNode block.IBlock) error {
	if topNode.GetHeight() < bc.GetHeight() {
		return errors.New("BlockChain the topnode is not height than current chain")
	}
	return nil
}

func GetChainHeight(bc *POAChain) uint32 {
	return bc.GetHeight()
}



