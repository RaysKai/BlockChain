package block

import (
	"github.com/linkchain/meta/tx"
	"github.com/linkchain/common/serialize"
)

type ethBlock struct{
}


func New()(IBlock, error){
	block := &ethBlock{}
	return block, nil
}

func (*ethBlock)SetTx([]tx.ITransaction)(error){
	return nil
}

func (b *ethBlock)GetHeight() int{
	return 0
}

func (*ethBlock)GetBlockID() IBlockID{
	return nil
}

func (*ethBlock)Verify()(error){
	return nil
}
//Serialize/Deserialize
func (*ethBlock)Serialize()(serialize.SerializeStream){
	return nil
}

func (*ethBlock)Deserialize(s serialize.SerializeStream){
}

//
func (*ethBlock)ToString()(string){
	return ""
}