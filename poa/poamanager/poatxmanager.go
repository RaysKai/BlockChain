package poamanager

import (
	"time"

	"github.com/linkchain/meta/tx"
	"github.com/linkchain/common/util/log"
	poameta "github.com/linkchain/poa/meta"
	"github.com/linkchain/meta/account"
	"github.com/linkchain/meta"
	"github.com/linkchain/common/math"
)

type POATxManager struct {
	txpool []poameta.POATransaction
}

/** interface: common.IService **/
func (m *POATxManager) Init(i interface{}) bool{
	log.Info("POABlockManager init...");
	m.txpool = make([]poameta.POATransaction,0)

	return true
}

func (m *POATxManager) Start() bool{
	log.Info("POABlockManager start...");
	return true
}

func (m *POATxManager) Stop(){
	log.Info("POABlockManager stop...");
}

func (m *POATxManager) AddTransaction(tx tx.ITx) error{
	newTx := *tx.(*poameta.POATransaction)
	m.txpool = append(m.txpool,newTx)
	return nil
}


func (m *POATxManager) GetAllTransaction() []tx.ITx{
	txs := make([]tx.ITx,0)
	for _,tx := range m.txpool {
		txs = append(txs,&tx)
	}
	return txs
}


func (m *POATxManager) RemoveTransaction(txid tx.ITxID) error{
	txidHash := txid.(math.Hash)
	deleteIndex := make([]int,0)
	for index,tx := range m.txpool{
		txHash := tx.GetTxID().(math.Hash)
		if txHash.IsEqual(&txidHash){
			deleteIndex = append(deleteIndex,index)
		}
	}
	for _,index := range deleteIndex{
		m.txpool = append(m.txpool[:index],m.txpool[index+1:]...)
	}
	return nil
}

func (m *POATxManager) NewTransaction(form account.IAccount,to account.IAccount,amount meta.IAmount) tx.ITx {
	newTx := poameta.POATransaction{Version:0,
		From:poameta.GetPOATransactionPeer(form,nil),
		To:poameta.GetPOATransactionPeer(to,nil),
		Amount:*amount.(*poameta.POAAmount),
		Time:time.Now()}
	return &newTx
}

func (m *POATxManager) CheckTx(tx tx.ITx) bool {
	log.Info("poa CheckTx ...")
	return true
}

func (m *POATxManager) ProcessTx(tx tx.ITx) {
	log.Info("poa ProcessTx ...")
	//1.checkTx
	if !m.CheckTx(tx) {
		log.Error("poa checkTransaction failed")
		return
	}
	//2.push Tx into storage
	m.AddTransaction(tx)
}

func (m *POATxManager) SignTransaction(tx tx.ITx) error  {
	return nil
}