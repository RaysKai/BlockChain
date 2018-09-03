package consensus

import (
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/common"
	"github.com/linkchain/poa"
	"github.com/linkchain/meta/block"
	"github.com/linkchain/meta/tx"
)

var (
	service poa.Service
)

type Service struct{

}

type ConsensusService interface {
	common.IService
	ProcessBlock(block block.IBlock)
	ProcessTx(tx tx.ITx)
}

func (s *Service) Init(i interface{}) bool{
	//log.Info("consensus service init...");
	service = poa.Service{}
	service.Init(i)
	return true
}

func (s *Service) Start() bool{
	//log.Info("consensus service start...");
	service.Start()
	return true
}

func (s *Service) Stop(){
	//log.Info("consensus service stop...");
	service.Stop()
}

func (s *Service) ProcessBlock(block block.IBlock){
	log.Info("ProcessBlock ...");
	service.ProcessBlock(block)
}

func (s *Service) ProcessTx(tx tx.ITx){
	log.Info("ProcessTx ...");
	service.ProcessTx(tx)
}

