package consensus

import (
	"github.com/linkchain/common/util/log"
)

type Service struct{
}

func (s *Service) Init(i interface{}) bool{
	log.Info("consensus service init...");
	return true
}

func (s *Service) Start() bool{
	log.Info("consensus service start...");
	return true
}

func (s *Service) Stop(){
	log.Info("consensus service stop...");
}