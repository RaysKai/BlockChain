package poamanager

import (
	"github.com/linkchain/common/util/log"
	poameta "github.com/linkchain/poa/meta"
	"github.com/linkchain/common/math"
	"github.com/linkchain/meta/block"
)

type POAChainManager struct {
	chains []poameta.POAChain
}

func (m *POAChainManager) Init(i interface{}) bool{
	log.Info("POAChainManager init...");
	m.chains = make([]poameta.POAChain,0)
	//create gensis chain
	gensisChain := poameta.NewPOAChain(GetManager().BlockMgr.GetGensisBlock(),GetManager().BlockMgr.GetGensisBlock())
	m.chains = append(m.chains,gensisChain)
	//todo need to load storage

	return true
}

func (m *POAChainManager) Start() bool{
	log.Info("POAChainManager start...");
	return true
}

func (s *POAChainManager) Stop(){
	log.Info("POAChainManager stop...");
}

func (m *POAChainManager) GetBestBlock() block.IBlock  {
	bestHeight := uint32(0);
	var bestBlock block.IBlock
	for _,chain := range m.chains {
		if bestHeight <= chain.GetHeight() {
			bestHeight = chain.GetHeight()
			bestBlock = chain.GetLastBlock()
		}
	}
	return bestBlock
}


func (m *POAChainManager) GetBlockByHash(h math.Hash) block.IBlock  {

	return nil
}

func (m *POAChainManager) GetBlockByHeight(height uint32) block.IBlock  {

	return nil
}

func (m *POAChainManager) GetBlockChainInfo() string  {
	return "this is poa chain";
}

func (m *POAChainManager) AddBlock(block block.IBlock)  {

}

func (m *POAChainManager) UpdateChain() bool  {
	return true
}
