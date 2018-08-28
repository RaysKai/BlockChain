package poa

import (
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/consensus/validator"
	"github.com/linkchain/meta"
	"github.com/linkchain/consensus/poa/poamanager"
)

type Service struct{
	blockValidator validator.BlockValidator
	transactionValidator validator.TransactionVlidator
}

func (s *Service) Init(i interface{}) bool{
	log.Info("poa consensus service init...")
	s.blockValidator = &POAValidator{}
	s.transactionValidator = &POAValidator{}
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

func (s *Service) ProcessBlock(block interface{}){
	log.Info("poa ProcessBlock ...")
	//1.checkBlock
	if !s.blockValidator.CheckBlock(block.(meta.Block)) {
		log.Error("poa checkBlock failed")
		return
	}
	//2.updateChain
	if poamanager.GetInstance().ChainMgr.UpdateChain() {
		log.Info("update chain successed")
	}
	//3.updateStorage
}

func (s *Service) ProcessTx(tx interface{}){
	log.Info("poa ProcessTx ...")
	//1.checkTx
	if !s.transactionValidator.CheckTx(tx.(meta.Transaction)) {
		log.Error("poa checkTransaction failed")
		return
	}
	//2.push Tx into storage
}