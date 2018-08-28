package account

import (
	"github.com/linkchain/meta"
	"github.com/linkchain/common/serialize"
)

type IAccount interface {
	//acount management
	ChangeAmount(meta.IAmount) meta.IAmount
	GetAmount() meta.IAmount

	//verifiy
	Verify()(error)

	//serialize
	serialize.ISerialize
}

