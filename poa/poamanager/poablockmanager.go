package poamanager

import (
	"sync"
	"time"
	"errors"
	"crypto/sha256"

	"github.com/linkchain/meta/block"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/common/math"
	poameta "github.com/linkchain/poa/meta"
)

const (
	MaxMapSize = 1024 * 4
)

type POABlockManager struct {
	sync.RWMutex
	mapBlockIndexByHash map[math.Hash]poameta.POABlock
}

func (m *POABlockManager) readMap(key math.Hash) (poameta.POABlock,bool) {
	m.RLock()
	value, ok := m.mapBlockIndexByHash[key]
	m.RUnlock()
	return value, ok
}

func (m *POABlockManager) writeMap(key math.Hash, value poameta.POABlock) {
	m.Lock()
	m.mapBlockIndexByHash[key] = value
	m.Unlock()
}


/** interface: common.IService **/
func (m *POABlockManager) Init(i interface{}) bool{
	log.Info("POABlockManager init...");
	m.mapBlockIndexByHash = make(map[math.Hash]poameta.POABlock)
	//load gensis block
	gensisBlock := GetManager().BlockManager.GetGensisBlock()
	m.AddBlock(gensisBlock)
	//load block by chainmanager

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
	bestBlock := GetManager().ChainManager.GetBestBlock()
	if bestBlock != nil {
		bestHash := bestBlock.GetBlockID().(math.Hash)
		txs := []poameta.POATransaction{}
		block := &poameta.POABlock{
			Header: poameta.POABlockHeader{Version: 0, PrevBlock: bestHash, MerkleRoot: math.Hash{}, Timestamp: time.Now(), Difficulty: 0x207fffff, Nonce: 0, Extra: nil, Height: bestBlock.GetHeight() + 1},
			TXs:    txs,
		}
		return block
	}else {
		return m.GetGensisBlock()
	}

}

/** interface: BlockBaseManager **/
func (m *POABlockManager) GetGensisBlock() block.IBlock{
	txs := []poameta.POATransaction{}
	fromAddress := math.Hash(sha256.Sum256([]byte("lf")))
	toAddress := math.Hash(sha256.Sum256([]byte("lc")))
	formAccount := &poameta.POAAccount{AccountID:poameta.POAAccountID{ID:fromAddress}}
	toAccount := &poameta.POAAccount{AccountID:poameta.POAAccountID{ID:toAddress}}
	amount := &poameta.POAAmount{Value:10}
	tx := *GetManager().TransactionManager.NewTransaction(formAccount,toAccount,amount).(*poameta.POATransaction)
	txs = append(txs, tx)
	block := &poameta.POABlock{
		Header:poameta.POABlockHeader{Version:0, PrevBlock:math.Hash{},MerkleRoot:math.Hash{},Timestamp:time.Now(),Difficulty:0x207fffff,Nonce:0,Extra:[]byte(string("hello,I am gensis block"))},
		TXs:txs,
	}
	return block
}

/** interface: BlockPoolManager **/
func (m *POABlockManager) GetBlockByID(hash block.IBlockID) (block.IBlock,error) {
	index, ok := m.readMap(hash.(math.Hash))
	if ok {
		return &index,nil
	}

	//TODO need to storage
	return nil, errors.New("POABlockManager can not find block by hash:" + hash.GetString())
}

func (m *POABlockManager) GetBlockByHeight(height uint32) ([]block.IBlock,error) {
	//TODO may not be need
	return nil,nil
}


func (m *POABlockManager) AddBlock(block block.IBlock) error{
	hash := block.GetBlockID().(math.Hash)
	m.writeMap(hash,*(block.(*poameta.POABlock)))
	return nil
}

func (m *POABlockManager) AddBlocks(blocks []block.IBlock) error{
	for _,block := range blocks {
		m.AddBlock(block)
	}
	return nil
}


func (m *POABlockManager) RemoveBlock(hash block.IBlockID) error{
	//TODO need to lock
	m.Lock()
	delete(m.mapBlockIndexByHash,*(hash.(*math.Hash)))
	m.Unlock()
	return nil
}

/** interface: BlockValidator **/
func (m *POABlockManager) CheckBlock(block block.IBlock) bool {
	log.Info("POA CheckBlock ...")

	return true
}

func (s *POABlockManager) ProcessBlock(block block.IBlock){
	log.Info("POA ProcessBlock ...")
	//1.checkBlock
	if !GetManager().BlockManager.CheckBlock(block) {
		log.Error("POA checkBlock failed")
		return
	}

	//2.acceptBlock
	GetManager().ChainManager.AddBlock(block)
	log.Info("POA Add a Blocks","block hash",block.GetBlockID().GetString())
	log.Info("POA Add a Blocks","prev hash",block.GetPrevBlockID().GetString())

	//3.updateChain
	if !GetManager().ChainManager.UpdateChain() {
		log.Info("POA Update chain failed")
		return
	}
	log.Info("POA ProcessBlock successed","blockchaininfo", GetManager().ChainManager.GetBlockChainInfo())

	//4.updateStorage
}