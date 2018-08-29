package poa

import (
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/poa/poamanager"
	"github.com/linkchain/meta/block"
	"github.com/linkchain/meta/tx"
)

type Service struct{

}

func (s *Service) Init(i interface{}) bool{
	log.Info("poa consensus service init...")
	poamanager.GetInstance().Init(nil)
	return true
}

func (s *Service) Start() bool{
	log.Info("poa consensus service start...")
	poamanager.GetInstance().Start()
	return true
}

func (s *Service) Stop(){
	log.Info("poa consensus service stop...")
	poamanager.GetInstance().Stop()
}

func (s *Service) ProcessBlock(block block.IBlock){
	log.Info("poa ProcessBlock ...")
	//1.checkBlock
	if !poamanager.GetInstance().BlockMgr.CheckBlock(block) {
		log.Error("poa checkBlock failed")
		return
	}
	//2.updateChain
	if poamanager.GetInstance().ChainMgr.UpdateChain() {
		log.Info("update chain successed")
	}
	//3.updateStorage
}

func (s *Service) ProcessTx(tx tx.ITransaction){
	log.Info("poa ProcessTx ...")
	//1.checkTx
	if !poamanager.GetInstance().TransactionMgr.CheckTx(tx) {
		log.Error("poa checkTransaction failed")
		return
	}
	//2.push Tx into storage
}