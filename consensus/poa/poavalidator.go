package poa

import (
	"github.com/linkchain/meta"
	"github.com/linkchain/common/util/log"
)

type POAValidator struct {

}

func (validator *POAValidator) CheckBlock(block meta.Block) bool {
	log.Info("poa checkBlock ...");
	log.Info(block.ToString())
	return true
}

func (validator *POAValidator) CheckTx(tx meta.Transaction) bool {
	log.Info("poa CheckTx ...");
	return true
}