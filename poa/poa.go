package poa

import (
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/poa/poamanager"
)

type Service struct{

}

func (s *Service) Init(i interface{}) bool{
	log.Info("poa consensus service init...")
	poamanager.GetManager().Init(nil)
	return true
}

func (s *Service) Start() bool{
	log.Info("poa consensus service start...")
	poamanager.GetManager().Start()
	return true
}

func (s *Service) Stop(){
	log.Info("poa consensus service stop...")
	poamanager.GetManager().Stop()
}

func (s *Service) GetManager() *poamanager.POAManager {
	return poamanager.GetManager()
}