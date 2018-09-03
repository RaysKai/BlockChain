package node

import (

	"github.com/linkchain/common"
	"github.com/linkchain/p2p"
	"github.com/linkchain/consensus"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/poa/poamanager"
)


var (
	//service collection
	 svcList = []common.IService{
		 &p2p.Service{},
		 &consensus.Service{},
	 }
)


func Init() {
	log.Info("Node init...")

	//init all service
	for _, v := range svcList {
		v.Init(nil)
	}
}

func Run(){
	log.Info("Node is running...")

	//start all service
	for _,v := range svcList{
		v.Start()
	}
	block :=poamanager.GetManager().BlockMgr.NewBlock()
	poamanager.GetManager().BlockMgr.ProcessBlock(block)
}
