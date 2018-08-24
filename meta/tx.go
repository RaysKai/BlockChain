package meta


type Transaction struct {
	// Version of the Transaction.  This is not the same as the Block version.
	Version int32

	From []TransactionIn

	To []TransactionOut

	// Extra used to extenion the block.
	Extra []byte

	Signs []FromSign
}

type TransactionIn struct {
	Id Hash
	Extra []byte
}

type TransactionOut struct {
	Address Hash
	Value uint32
	Extra []byte
}

type FromSign struct {
	Code []byte
}