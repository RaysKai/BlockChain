package poamanager

import (
	"time"

	"github.com/linkchain/meta/block"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/common/math"
	poameta "github.com/linkchain/poa/meta"
	"errors"
)

var mapBlockIndexByHash map[math.Hash]int
//var mapBlockIndexByHeight map[meta.Hash]int
var blockIndexs []poameta.BlockIndex

type POABlockManager struct {

}

func (m *POABlockManager) Init(i interface{}) bool{
	log.Info("POABlockManager init...");
	mapBlockIndexByHash = make(map[math.Hash]int)
	blockIndexs = make([]poameta.BlockIndex,0)
	return true
}

func (m *POABlockManager) Start() bool{
	log.Info("POABlockManager start...");
	return true
}

func (m *POABlockManager) Stop(){
	log.Info("POABlockManager stop...");
}


func (m *POABlockManager) GetBlock(hash block.IBlockID) (*poameta.BlockIndex,error) {
	index, ok := mapBlockIndexByHash[(*hash.(*math.Hash))]
	if ok {
		return &blockIndexs[index],nil
	}
	return &blockIndexs[index], errors.New("POABlockManager can not find block by hash:" + hash.GetString())
}

func (m *POABlockManager) AddBlock(block block.IBlock){
	index := len(blockIndexs)
	blockIndex := poameta.NewBlockIndex(*block.(*poameta.POABlock),poameta.StateLess)
	blockIndexs = append(blockIndexs,blockIndex)
	mapBlockIndexByHash[(*block.GetBlockID().(*math.Hash))] = index
}

func (m *POABlockManager) RemoveBlock(hash block.IBlockID){

}


func (m *POABlockManager) CheckBlock(block block.IBlock) bool {
	log.Info("poa CheckBlock ...")
	GetInstance().ChainMgr.AddBlock(block)
	log.Info("poa Add a Block","block",block)
	log.Info("poa GetBestHeader","blockhash", GetInstance().ChainMgr.GetBestBlock().GetBlockID().GetString())
	txs := []poameta.POATransaction{}
	block1 := &poameta.POABlock{
		Header:poameta.POABlockHeader{Version:1, PrevBlock:math.Hash{},MerkleRoot:math.Hash{},Timestamp:time.Unix(1401292357, 0),Difficulty:0x207fffff,Nonce:0,Extra:nil},
		TXs:txs,
	}
	log.Info("poa Add a Block","block",block1)
	GetInstance().ChainMgr.AddBlock(block1)
	log.Info("poa GetBestHeader","blockhash",GetInstance().ChainMgr.GetBestBlock().GetBlockID().GetString())
	return true
}

func (s *POABlockManager) ProcessBlock(block block.IBlock){
	log.Info("poa ProcessBlock ...")
	//1.checkBlock
	if !GetInstance().BlockMgr.CheckBlock(block) {
		log.Error("poa checkBlock failed")
		return
	}
	//2.updateChain
	if GetInstance().ChainMgr.UpdateChain() {
		log.Info("update chain successed")
	}
	//3.updateStorage
}