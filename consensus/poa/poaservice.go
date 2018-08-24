package poa

import (
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/consensus/validator"
	"github.com/linkchain/consensus/manager"
	"github.com/linkchain/meta"
)

type Service struct{
	blockValidator validator.BlockValidator
	transactionValidator validator.TransactionVlidator
	blockManager manager.BlockManager
	accountManager manager.AccountManager
	transactionManager manager.TransactionManager
	chainManager manager.ChainManager
}

func (s *Service) Init(i interface{}) bool{
	log.Info("poa consensus service init...");
	s.blockValidator = &POAValidator{};
	s.transactionValidator = &POAValidator{};
	return true
}

func (s *Service) Start() bool{
	log.Info("poa consensus service start...");
	return true
}

func (s *Service) Stop(){
	log.Info("poa consensus service stop...");
}

func (s *Service) ProcessBlock(block meta.Block){
	log.Info("poa ProcessBlock ...");
	//1.checkBlock
	if !s.blockValidator.CheckBlock(block) {
		log.Error("poa checkBlock failed")
		return
	}
	//2.updateChain
	//3.updateStorage
}

func (s *Service) ProcessTx(tx meta.Transaction){
	log.Info("poa ProcessTx ...");
	//1.checkTx
	if !s.transactionValidator.CheckTx(tx) {
		log.Error("poa checkTransaction failed")
		return
	}
	//2.push Tx into storage
}