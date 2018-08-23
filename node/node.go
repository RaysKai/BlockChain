package node

import (
	"github.com/linkchain/common"
	"github.com/linkchain/p2p"
	"github.com/linkchain/consensus"
	"github.com/linkchain/common/util/log"
)

var (
	svcList []common.IService
)


func Init(){
	log.Info("Node init...")

	svcList = append(svcList, &p2p.Service{})
	svcList = append(svcList, &consensus.Service{})

	for _,v := range svcList{
		v.Init(nil)
	}
}

func Run(){
	log.Info("Node is running...")
	for _,v := range svcList{
		v.Start()
	}
}
