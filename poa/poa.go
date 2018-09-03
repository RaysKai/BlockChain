package poa

import (
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/poa/poamanager"
)

type Service struct{

}

func (s *Service) Init(i interface{}) bool{
	log.Info("poa consensus service init...")
	s.GetManager().Init(nil)
	return true
}

func (s *Service) Start() bool{
	log.Info("poa consensus service start...")
	s.GetManager().Start()
	return true
}

func (s *Service) Stop(){
	log.Info("poa consensus service stop...")
	s.GetManager().Stop()
}

func (s *Service) GetManager() *poamanager.POAManager {
	return poamanager.GetManager()
}