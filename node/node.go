package node

import (
	"fmt"
	"github.com/linkchain/common"
	"github.com/linkchain/p2p"
	"github.com/linkchain/consensus"
)

var (
	svcList []common.IService
)


func Init(){
	fmt.Println("Node init...")

	svcList = append(svcList, &p2p.Service{})
	svcList = append(svcList, &consensus.Service{})

	for _,v := range svcList{
		v.Init(nil)
	}
}

func Run(){
	fmt.Println("Node is running...")
	for _,v := range svcList{
		v.Start()
	}
}
