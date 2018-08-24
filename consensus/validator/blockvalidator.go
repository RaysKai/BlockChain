package validator

import "github.com/linkchain/meta"

type BlockValidator interface {
	CheckBlock(block meta.Block) bool
}

