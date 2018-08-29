package validator

import "github.com/linkchain/meta/block"

type BlockValidator interface {
	CheckBlock(block block.IBlock) bool
}

