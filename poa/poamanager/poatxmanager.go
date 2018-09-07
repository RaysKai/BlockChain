package poamanager

import (
	"github.com/linkchain/meta/tx"
	"github.com/linkchain/common/util/log"
	poameta "github.com/linkchain/poa/meta"
	"github.com/linkchain/meta/account"
	"github.com/linkchain/meta"
)

type POATxManager struct {
}


func (m *POATxManager) NewTransaction(form account.IAccount,to account.IAccount,amount meta.IAmount) tx.ITx {
	newTx := poameta.POATransaction{Version:0,
		From:poameta.GetPOATransactionPeer(form,nil),
		To:poameta.GetPOATransactionPeer(to,nil),
		Amount:*amount.(*poameta.POAAmount)}
	return &newTx
}

func (m *POATxManager) CheckTx(tx tx.ITx) bool {
	log.Info("poa CheckTx ...")
	return true
}

func (m *POATxManager) ProcessTx(tx tx.ITx){
	log.Info("poa ProcessTx ...")
	//1.checkTx
	if !GetManager().TransactionManager.CheckTx(tx) {
		log.Error("poa checkTransaction failed")
		return
	}
	//2.push Tx into storage
}

func (m *POATxManager) SignTransaction(tx tx.ITx) error  {
	return nil
}