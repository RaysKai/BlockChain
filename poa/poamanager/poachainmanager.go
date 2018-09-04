package poamanager

import (
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/common/math"
	"github.com/linkchain/meta/block"
	poameta "github.com/linkchain/poa/meta"
)

type POAChainManager struct {
	chains []poameta.POAChain
	mainChain []poameta.POAChainNode
}

func (m *POAChainManager) Init(i interface{}) bool{
	log.Info("POAChainManager init...");

	//create gensis chain
	gensisBlock := GetManager().BlockManager.GetGensisBlock()

	gensisChain := poameta.NewPOAChain(gensisBlock,gensisBlock)
	m.chains = make([]poameta.POAChain,0)
	m.chains = append(m.chains,gensisChain)

	gensisChainNode := poameta.NewPOAChainNode(*gensisBlock.(*poameta.POABlock))
	m.mainChain = make([]poameta.POAChainNode,0)
	m.mainChain = append(m.mainChain,gensisChainNode)

	//TODO need to load storage
	m.UpdateChain()
	return true
}

func (m *POAChainManager) Start() bool{
	log.Info("POAChainManager start...");
	//TODO need to updateMainChain
	return true
}

func (s *POAChainManager) Stop(){
	log.Info("POAChainManager stop...");
}

func (m *POAChainManager) GetBestBlock() block.IBlock  {
	hash := m.mainChain[m.GetBestHeight()].GetNodeHash()
	bestBlock,_ := GetManager().BlockManager.GetBlockByID(hash)
	return bestBlock
}

func (m *POAChainManager) GetBestBlockHash() block.IBlockID  {
	hash := m.mainChain[m.GetBestHeight()].GetNodeHash()
	return hash
}

func (m *POAChainManager) GetBestHeight() uint32 {
	return uint32(len(m.mainChain)-1)
}

func (m *POAChainManager) GetBlockByHash(hash math.Hash) block.IBlock  {
	block,_ := GetManager().BlockManager.GetBlockByID(hash)
	return block
}

func (m *POAChainManager) GetBlockByHeight(height uint32) block.IBlock  {
	block,_ := GetManager().BlockManager.GetBlockByID(m.mainChain[height].CurentHash)
	return block
}

func (m *POAChainManager) GetBlockChainInfo() string  {
	return "this is poa chain";
}

func (m *POAChainManager) AddBlock(block block.IBlock)  {
	poablock := *block.(*poameta.POABlock)
	blockPrevHash := block.GetPrevBlockID().(math.Hash)
	GetManager().BlockManager.AddBlock(&poablock)
	bestBlockHash := m.GetBestBlockHash().(math.Hash)
	//This Block is on mainchain
	if bestBlockHash.IsEqual(&blockPrevHash) {
		m.mainChain = append(m.mainChain,poameta.NewPOAChainNode(poablock))

		_,mainChainIndex := m.GetLongestChain()
		m.chains[mainChainIndex].Block = append(m.chains[mainChainIndex].Block,poablock)
		return
	}
	//This Block is on SideChain
	//TODO 1.Find a sideChain where the block is added.if block is no_parent maybe cache
	//TODO 2.update chains and push the block on sidechain
	//TODO 3.updaeChain
}

func (m *POAChainManager) GetLongestChain() (poameta.POAChain,int)  {
	var mainChain poameta.POAChain
	bestHeight := uint32(0);
	position := 0
	for i,chain := range m.chains {
		if bestHeight <= chain.GetHeight() {
			bestHeight = chain.GetHeight()
			mainChain = chain
			position = i
		}
	}
	return mainChain,position
}

func (m *POAChainManager) UpdateChain() bool  {
	chain,_ := m.GetLongestChain()
	block := chain.Block[len(chain.Block)-1]
	bestNode := poameta.NewPOAChainNode(block)

	//1.check best Height
	if (bestNode.Height+1) > uint32(len(m.mainChain)) {
		//update chain
		m.AddBlock(&block)
		return m.sortChain()
	}
	hash := m.GetBestBlockHash().(math.Hash)
	//2.check best Hash
	if !bestNode.CurentHash.IsEqual(&hash) {
		return m.sortChain()
	}
	return true
}

func (m *POAChainManager) sortChain() bool  {
	currentHeight := m.GetBestHeight()
	currentNode := m.mainChain[currentHeight]
	prevNode := m.mainChain[currentHeight-1]

	for currentHeight >= 0 {
		if currentNode.CheckPrev(prevNode) {
			break
		}
		realPrevBlock,error := GetManager().BlockManager.GetBlockByID(currentNode.PrevHash)
		if error != nil {
			log.Error("POA sort Chain error,you need to download data")
			return false
		}
		m.mainChain[currentHeight-1] = poameta.NewPOAChainNode(*realPrevBlock.(*poameta.POABlock))
		currentHeight--
	}
	return true
}


