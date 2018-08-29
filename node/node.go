package node

import (
	"time"

	"github.com/linkchain/common"
	"github.com/linkchain/p2p"
	"github.com/linkchain/consensus"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/common/math"
	poameta "github.com/linkchain/poa/meta"
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
	txs := []poameta.POATransaction{}
	block := &poameta.POABlock{
		Header:poameta.POABlockHeader{Version:0, PrevBlock:math.Hash{},MerkleRoot:math.Hash{},Timestamp:time.Unix(1401292357, 0),Difficulty:0x207fffff,Nonce:0,Extra:nil},
		TXs:txs,
	}
	//log.Info(block.ToString())
	svcList[1].(consensus.ConsensusService).ProcessBlock(block)
}
