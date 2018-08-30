package account

import (
	"github.com/linkchain/meta"
	"github.com/linkchain/common/serialize"
)

type IAccountID interface{
	GetString() string
}
type IAccount interface {
	//acount management
	ChangeAmount(meta.IAmount) meta.IAmount
	GetAmount() meta.IAmount

	GetAccountID()  IAccountID
	//verifiy
	Verify()(error)

	//serialize
	serialize.ISerialize
}

