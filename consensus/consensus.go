package consensus

import (
	"github.com/linkchain/poa"
)

var (
	service poa.Service
)

type Service struct{

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

