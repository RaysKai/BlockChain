package meta

import "github.com/linkchain/common/math"

type POATransaction struct {
	// Version of the Transaction.  This is not the same as the Blocks version.
	Version int32

	From []POATransactionIn

	To []POATransactionOut

	// Extra used to extenion the block.
	Extra []byte

	Signs []FromSign
}

type POATransactionIn struct {
	Id math.Hash
	Extra []byte
}

type POATransactionOut struct {
	Address math.Hash
	Value uint32
	Extra []byte
}

type FromSign struct {
	Code []byte
}
