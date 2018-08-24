package meta

import (
	"encoding/json"
)

type Block struct {
	Header BlockHeader
	TXs []Transaction
}

func (block *Block) ToString() string  {
	data, err := json.Marshal(block);
	if  err != nil {
		return err.Error()
	}
	return string(data)
}