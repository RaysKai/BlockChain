package block

import (
	"github.com/linkchain/meta/tx"
	"github.com/linkchain/common/serialize"
)

type IBlockID interface{
	GetString() string
}

type IBlock interface {
	//block content
	SetTx([]tx.ITx)(error)

	GetHeight() uint32

	GetBlockID() IBlockID

	GetPrevBlockID() IBlockID
	//verifiy
	Verify()(error)

	//serialize
	serialize.ISerialize
}