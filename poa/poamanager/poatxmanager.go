package poamanager

import (
	"github.com/linkchain/meta/tx"
	"github.com/linkchain/common/util/log"
)

type POATxManager struct {
}


func (m *POATxManager) NewTransaction() tx.ITx {
	return nil
}

func (m *POATxManager) CheckTx(tx tx.ITx) bool {
	log.Info("poa CheckTx ...")
	return true
}

func (m *POATxManager) ProcessTx(tx tx.ITx){
	log.Info("poa ProcessTx ...")
	//1.checkTx
	if !GetInstance().TransactionMgr.CheckTx(tx) {
		log.Error("poa checkTransaction failed")
		return
	}
	//2.push Tx into storage
}

func (m *POATxManager) SignTransaction(tx tx.ITx) error  {
	return nil
}