package meta

import (
	"encoding/json"
	"bytes"
	"encoding/binary"
	"crypto/sha256"
	"time"

	"github.com/linkchain/meta/tx"
	"github.com/linkchain/common/serialize"
	"github.com/linkchain/common/math"
	"github.com/linkchain/meta/block"
)

type POABlock struct{
	Header POABlockHeader
	TXs []POATransaction
}

type POABlockHeader struct {
	// Version of the block.  This is not the same as the protocol version.
	Version int32

	// Hash of the previous block header in the block chain.
	PrevBlock math.Hash

	// Merkle tree reference to hash of all transactions for the block.
	MerkleRoot math.Hash

	// Time the block was created.  This is, unfortunately, encoded as a
	// uint32 on the wire and therefore is limited to 2106.
	Timestamp time.Time

	// Difficulty target for the block.
	Difficulty uint32

	// Nonce used to generate the block.
	Nonce uint32

	// Extra used to extenion the block.
	Extra []byte
}


func New()(block.IBlock, error){
	block := &POABlock{}
	return block, nil
}

func (b *POABlock)SetTx([]tx.ITransaction)(error){
	return nil
}

func (b *POABlock)GetHeight() int{
	return 0
}

func (b *POABlock)GetBlockID() block.IBlockID{
	data := make([]byte, 0)
	buf := bytes.NewBuffer(data)
	binary.Write(buf, binary.BigEndian, b.Header.Version)

	first := sha256.Sum256(data)
	return math.Hash(sha256.Sum256(first[:]))
}

func (b *POABlock)Verify()(error){
	return nil
}
//Serialize/Deserialize
func (b *POABlock)Serialize()(serialize.SerializeStream){
	return nil
}

func (b *POABlock)Deserialize(s serialize.SerializeStream){
}

//
func (b *POABlock)ToString()(string){
	data, err := json.Marshal(b);
	if  err != nil {
		return err.Error()
	}
	return string(data)
}