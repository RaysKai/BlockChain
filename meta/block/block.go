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
	SetTx([]tx.ITransaction)(error)

	GetHeight() int

	GetBlockID() IBlockID

	//verifiy
	Verify()(error)

	//serialize
	serialize.ISerialize
}