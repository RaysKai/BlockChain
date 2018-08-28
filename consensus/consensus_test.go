package consensus

import (
	"testing"
	"time"
	"github.com/linkchain/meta"
)

var testService Service = Service{}

// ProcessBlock test
func TestProcessBlock(t *testing.T)  {
	txs := []meta.Transaction{}
	block := meta.Block{
		Header:meta.BlockHeader{Version:0, PrevBlock:meta.Hash{},MerkleRoot:meta.Hash{},Timestamp:time.Unix(1401292357, 0),Difficulty:0x207fffff,Nonce:0,Extra:nil},
		TXs:txs,
	}
	testService.Init(nil)
	testService.Start()
	testService.ProcessBlock(block)
	t.Log(block.ToString())
	testService.Stop()
}
