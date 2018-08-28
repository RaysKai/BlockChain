package consensus

import (
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/consensus/poa"
	"github.com/linkchain/common"
)

var (
	service poa.Service
)

type Service struct{

}

type ConsensusService interface {
	common.IService
	ProcessBlock(block interface{})
	ProcessTx(tx interface{})
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

func (s *Service) ProcessBlock(block interface{}){
	log.Info("ProcessBlock ...");
	service.ProcessBlock(block)
}

func (s *Service) ProcessTx(tx interface{}){
	log.Info("ProcessTx ...");
	service.ProcessTx(tx)
}

