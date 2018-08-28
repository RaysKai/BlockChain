package poamanager

import (
	"github.com/linkchain/meta"
	"github.com/linkchain/common/util/log"
)

var mainchain  []meta.Block
var bestHeight int


type POAChainManager struct {

}

func (m *POAChainManager) Init(i interface{}) bool{
	log.Info("POAChainManager init...");
	mainchain = make([]meta.Block,1)
	bestHeight=0;
	return true
}

func (m *POAChainManager) Start() bool{
	log.Info("POAChainManager start...");
	return true
}

func (s *POAChainManager) Stop(){
	log.Info("POAChainManager stop...");
}

func (m *POAChainManager) GetBestHeader() *meta.BlockHeader  {
	return &mainchain[bestHeight].Header;
}

func (m *POAChainManager) GetMainChain() *meta.BlockHeader  {
	return &mainchain[bestHeight].Header;
}

func (m *POAChainManager) GetBlockByHash(h *meta.Hash) meta.Block  {
	return mainchain[bestHeight];
}

func (m *POAChainManager) GetBlockByHeight(height uint32) meta.Block  {
	return mainchain[height];
}

func (m *POAChainManager) GetBlockChainInfo() string  {
	return "this is poa chain";
}

func (m *POAChainManager) AddBlock(block meta.Block)  {
	mainchain = append(mainchain,block)
	bestHeight++
}

func (m *POAChainManager) UpdateChain() bool  {
	return true
}
