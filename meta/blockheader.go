package meta

import (
	"time"
	"encoding/json"
)

type BlockHeader struct {
	// Version of the block.  This is not the same as the protocol version.
	Version int32

	// Hash of the previous block header in the block chain.
	PrevBlock Hash

	// Merkle tree reference to hash of all transactions for the block.
	MerkleRoot Hash

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

func (header *BlockHeader) ToString() string  {
	data, err := json.Marshal(header);
	if  err != nil {
		return err.Error()
	}
	return string(data)
}