package consensus

import (
	"github.com/linkchain/poa"
	"github.com/linkchain/consensus/manager"
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

func (s *Service) GetBlockManager() manager.BlockManager {
	return service.GetManager().BlockMgr
}

func (s *Service) GetTXManager() manager.TransactionManager {
	return service.GetManager().TransactionMgr
}

func (s *Service) GetAccountManager() manager.AccountManager {
	return service.GetManager().AccountMgr
}

func (s *Service) GetChainManager() manager.ChainManager {
	return service.GetManager().ChainMgr
}
