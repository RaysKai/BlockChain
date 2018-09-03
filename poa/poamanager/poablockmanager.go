package poamanager

import (
	"time"
	"errors"

	"github.com/linkchain/meta/block"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/common/math"
	poameta "github.com/linkchain/poa/meta"
)

var mapBlockIndexByHash map[math.Hash]uint32
var mapBlockIndexByHeight map[uint32][]uint32
//var blockIndexs []poameta.BlockIndex
var sideChain [][]uint32

type POABlockManager struct {

}

/** interface: common.IService **/
func (m *POABlockManager) Init(i interface{}) bool{
	log.Info("POABlockManager init...");
	mapBlockIndexByHash = make(map[math.Hash]uint32)
	//blockIndexs = make([]poameta.BlockIndex,0)
	mapBlockIndexByHeight = make(map[uint32][]uint32)
	sideChain = make([][]uint32,0,3)
	noParent := make([]uint32,0)//no-parent blocks | statusId 0x00000000
	mainChain := make([]uint32,0)//mainchain blocks | statusId 0x00000001
	fristSideChain := make([]uint32,0)//side(1) blocks | statusId 0x00000010
	sideChain = append(sideChain,noParent)
	sideChain = append(sideChain,mainChain)
	sideChain = append(sideChain,fristSideChain)
	return true
}

func (m *POABlockManager) Start() bool{
	log.Info("POABlockManager start...");
	return true
}

func (m *POABlockManager) Stop(){
	log.Info("POABlockManager stop...");
}

/** interface: BlockBaseManager **/
func (m *POABlockManager) NewBlock() block.IBlock{
	txs := []poameta.POATransaction{}
	block := &poameta.POABlock{
		Header:poameta.POABlockHeader{Version:0, PrevBlock:math.Hash{},MerkleRoot:math.Hash{},Timestamp:time.Now(),Difficulty:0x207fffff,Nonce:0,Extra:nil},
		TXs:txs,
	}
	return block
}

/** interface: BlockBaseManager **/
func (m *POABlockManager) GetGensisBlock() block.IBlock{
	txs := []poameta.POATransaction{}
	block := &poameta.POABlock{
		Header:poameta.POABlockHeader{Version:0, PrevBlock:math.Hash{},MerkleRoot:math.Hash{},Timestamp:time.Now(),Difficulty:0x207fffff,Nonce:0,Extra:[]byte(string("hello,I am gensis block"))},
		TXs:txs,
	}
	return block
}

/** interface: BlockPoolManager **/
func (m *POABlockManager) GetBlockByID(hash block.IBlockID) (block.IBlock,error) {
	/*index, ok := mapBlockIndexByHash[(*hash.(*math.Hash))]
	if ok {
		return nil,nil
	}*/
	return nil, errors.New("POABlockManager can not find block by hash:" + hash.GetString())
}

func (m *POABlockManager) GetBlockByHeight(height uint32) ([]block.IBlock,error) {
	/*index, ok := mapBlockIndexByHeight[height]
	if ok {
		var blocks []block.IBlock
		blocks = make([]block.IBlock,0)
		for i :=range index {
			blocks = append(blocks,&blockIndexs[i].Block)
		}
		if len(blocks) < 0 {
			return nil, errors.New("POABlockManager can not find block by height:" + string(height))
		}
		return blocks,nil
	}*/
	return nil, errors.New("POABlockManager can not find block by height:" + string(height))
}


func (m *POABlockManager) AddBlock(block block.IBlock) error{
	/*index := len(blockIndexs)
	status := poameta.StateLess
	//update maphash
	mapBlockIndexByHash[(*block.GetBlockID().(*math.Hash))] = uint32(index)
	//update mapheight
	mapBlockIndexByHeight[block.GetHeight()] = append(mapBlockIndexByHeight[block.GetHeight()],uint32(index))
	//update sidechain
	for _,index := range mapBlockIndexByHeight[(block.GetHeight()-1)] {
		if blockIndexs[index].Blockhash.IsEqual(block.GetPrevBlockID().(*math.Hash)) {
			status = blockIndexs[index].Status
			sideChain[blockIndexs[index].Status] = append(sideChain[blockIndexs[index].Status],index)
		}
	}
	if status == poameta.StateLess {

	}
	blockIndex := poameta.NewBlockIndex(*block.(*poameta.POABlock),status)
	blockIndexs = append(blockIndexs,blockIndex)*/
	return nil
}

func (m *POABlockManager) AddBlocks(block []block.IBlock) error{
	return nil
}


func (m *POABlockManager) RemoveBlock(hash block.IBlockID) error{
	return nil
}

/** interface: BlockValidator **/
func (m *POABlockManager) CheckBlock(block block.IBlock) bool {
	log.Info("POA CheckBlock ...")
	GetManager().ChainMgr.AddBlock(block)
	log.Info("POA Add a Block","block",block)
	log.Info("POA GetBestHeader","blockhash", GetManager().ChainMgr.GetBestBlock().GetBlockID().GetString())
	return true
}

func (s *POABlockManager) ProcessBlock(block block.IBlock){
	log.Info("POA ProcessBlock ...")
	//1.checkBlock
	if !GetManager().BlockMgr.CheckBlock(block) {
		log.Error("POA checkBlock failed")
		return
	}
	//2.updateChain
	if GetManager().ChainMgr.UpdateChain() {
		log.Info("Update chain successed")
	}
	//3.updateStorage
}