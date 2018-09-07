package meta

import (
	"github.com/linkchain/common/math"
	"github.com/linkchain/common/serialize"
	"github.com/linkchain/meta"
	"github.com/linkchain/meta/account"
	"encoding/json"
)


type POAAccountID struct {
	ID math.Hash
}

func (id *POAAccountID) GetString() string  {
	return id.ID.GetString()
}

type POAAccount struct {
	AccountID POAAccountID
	Value POAAmount
}

func (a *POAAccount) ChangeAmount(amount meta.IAmount) meta.IAmount{
	a.Value = *amount.(*POAAmount)
	return &a.Value
}

func (a *POAAccount) GetAmount() meta.IAmount{
	return &(a.Value)
}

func (a *POAAccount) GetAccountID() account.IAccountID{
	return &a.AccountID
}


//Serialize/Deserialize
func (a *POAAccount)Serialize()(serialize.SerializeStream){
	return nil
}

func (a *POAAccount)Deserialize(s serialize.SerializeStream){
}

func (a *POAAccount) ToString()(string) {
	data, err := json.Marshal(a);
	if  err != nil {
		return err.Error()
	}
	return string(data)
}


