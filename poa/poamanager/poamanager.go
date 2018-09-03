package poamanager

import (
	"sync"

	"github.com/linkchain/consensus/manager"
	"github.com/linkchain/common/util/log"
)

var m *POAManager
var once sync.Once

func GetManager() *POAManager {
	once.Do(func() {
		m = &POAManager {BlockMgr:&POABlockManager{},
			AccountMgr:&POAAccountManager{},
			TransactionMgr:&POATxManager{},
			ChainMgr:&POAChainManager{}}
	})
	return m
}

type POAManager struct {
	BlockMgr manager.BlockManager
	AccountMgr manager.AccountManager
	TransactionMgr manager.TransactionManager
	ChainMgr manager.ChainManager
}

func (m *POAManager) Init(i interface{}) bool{
	log.Info("POAManager init...");
	m.ChainMgr.Init(i)
	m.BlockMgr.Init(i)
	return true
}

func (m *POAManager) Start() bool{
	log.Info("POAManager start...");
	m.ChainMgr.Start()
	return true
}

func (m *POAManager) Stop(){
	log.Info("POAManager stop...");
	m.ChainMgr.Stop()
}

