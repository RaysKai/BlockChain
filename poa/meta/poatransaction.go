package meta

import (
	"github.com/linkchain/meta/account"
	"github.com/linkchain/meta"
	"github.com/linkchain/meta/tx"
	"github.com/linkchain/common/math"
	"github.com/linkchain/common/serialize"
	"crypto/sha256"
	"encoding/json"
)

type POATransactionPeer struct {
	AccountID POAAccountID
	Extra []byte
}

func GetPOATransactionPeer(iaccount account.IAccount, extra []byte) POATransactionPeer {
	id := *iaccount.GetAccountID().(*POAAccountID)
	return POATransactionPeer{AccountID:id,Extra:extra}
}

type FromSign struct {
	Code []byte
}

type POATransaction struct {
	// Version of the Transaction.  This is not the same as the Blocks version.
	Version int32

	From POATransactionPeer

	To POATransactionPeer

	Amount POAAmount

	// Extra used to extenion the block.
	Extra []byte

	Signs FromSign
}

func (tx *POATransaction) GetTxID() tx.ITxID  {
	first := sha256.Sum256(tx.To.AccountID.ID.CloneBytes())
	return math.Hash(sha256.Sum256(first[:]))
}

func (tx *POATransaction) SetFrom(from tx.ITxPeer)  {
	txin := *from.(*POATransactionPeer)
	tx.From = txin
}

func (tx *POATransaction) SetTo(to tx.ITxPeer)  {
	txout := *to.(*POATransactionPeer)
	tx.To = txout
}

func (tx *POATransaction) SetAmount(iAmount meta.IAmount)  {
	amount := *iAmount.(*POAAmount)
	tx.Amount = amount
}

func (tx *POATransaction) GetFrom() tx.ITxPeer  {
	return &tx.From
}

func (tx *POATransaction) GetTo() tx.ITxPeer  {
	return &tx.To
}

func (tx *POATransaction) GetAmount() meta.IAmount  {
	return &tx.Amount
}

func (tx *POATransaction) Sign()(math.ISignature, error)  {
	//TODO sign need to finish
	return nil,nil
}

func (tx *POATransaction) GetSignature()(math.ISignature)  {
	return nil
}

func (tx *POATransaction) Verify()(error)  {
	return nil
}

//Serialize/Deserialize
func (tx *POATransaction) Serialize()(serialize.SerializeStream){
	return nil
}

func (tx *POATransaction) Deserialize(s serialize.SerializeStream){
}

func (tx *POATransaction) ToString()(string) {
	data, err := json.Marshal(tx);
	if  err != nil {
		return err.Error()
	}
	return string(data)
}


