package meta

import (
	"errors"

	"github.com/linkchain/meta/block"
)

type POAChain struct {
	Blocks []POABlock
	IsInComplete bool
}

func NewPOAChain(startNode block.IBlock,endNode block.IBlock) POAChain {
	chainNode := make([]POABlock,0)
	isInComplete := false
	if startNode != nil {
		chainNode = append(chainNode, *(startNode.(*POABlock)))
		isInComplete = true
	}
	if endNode != nil {
		chainNode = append(chainNode, *(endNode.(*POABlock)))
	}
	return POAChain{Blocks:chainNode,IsInComplete:isInComplete}
}


func (bc *POAChain) AddNewBlock(block block.IBlock) {
	bc.Blocks = append(bc.Blocks,*(block.(*POABlock)))
}

/**invalidate block by block*/
func (bc *POAChain) Rollback(block.IBlock) {

}

/**invalidate block by height*/
func (bc *POAChain) RollbackAtHeight(int) {

}

func (bc *POAChain) GetLastBlock() block.IBlock {
	return &bc.Blocks[len(bc.Blocks)-1]
}

func (bc *POAChain) GetFirstBlock() block.IBlock {
	return &bc.Blocks[0]
}

func (bc *POAChain) GetHeight() uint32  {
	return bc.GetLastBlock().GetHeight()
}

func (bc *POAChain) GetBlockByID(block.IBlockID) block.IBlock {
	//TODO need to sorage
	return nil
}

func (bc *POAChain) GetBlockByHeight(int) block.IBlock {
	//TODO need to sorage
	return nil
}

func (bc *POAChain) UpdateChainTop(topBlock block.IBlock) error {
	if topBlock.GetHeight() < bc.GetHeight() {
		return errors.New("BlockChain the topnode is not height than current chain")
	}
	lastNode := NewPOAChainNode(bc.GetLastBlock())
	topNode := NewPOAChainNode(topBlock)
	if topNode.CheckPrev(lastNode) {
		bc.AddNewBlock(topBlock)
		return nil
	}else {
		return errors.New("BlockChain the topBlock is not next of lastBlock chain")
	}
}

func (bc *POAChain) AddChain(newChain POAChain) error {
	if bc.CanLink(newChain) {
		for _,block := range newChain.Blocks {
			bc.UpdateChainTop(&block)
		}
		return nil
	}else {
		return errors.New("BlockChain the topBlock is not next of lastBlock chain")
	}
}

func (bc *POAChain) CanLink(newChain POAChain) bool {
	if !newChain.IsInComplete{
		return false
	}
	topNode := NewPOAChainNode(&newChain.Blocks[1])
	lastNode := NewPOAChainNode(bc.GetLastBlock())
	return topNode.CheckPrev(lastNode)
}


/**
	GetNewChain
	get a new chain from this chain
 */
func (bc *POAChain) GetNewChain(forkBlock block.IBlock) POAChain{
	newChain := POAChain{IsInComplete:true}
	for _,block := range bc.Blocks {
		if block.GetHeight() < forkBlock.GetHeight() {
			newChain.AddNewBlock(&block)
		}
	}
	newChain.AddNewBlock(forkBlock)
	return newChain
}

func GetChainHeight(bc *POAChain) uint32 {
	return bc.GetHeight()
}



