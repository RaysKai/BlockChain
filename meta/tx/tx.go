package tx
import (
	"github.com/linkchain/common/serialize"
	"github.com/linkchain/meta"
	"github.com/linkchain/common/math"
)


type ITxPeer interface{

}

type ITxID interface{
	GetString() string
}

type ITx interface {

	GetTxID() ITxID

	//tx content
	SetFrom(from ITxPeer)
	SetTo(to ITxPeer)
	SetAmount(meta.IAmount)

	GetFrom() ITxPeer
	GetTo() ITxPeer
	GetAmount() meta.IAmount

	//signature
	Sign()(math.ISignature, error)
	GetSignature()(math.ISignature)
	Verify()(error)

	//serialize
	serialize.ISerialize
}
