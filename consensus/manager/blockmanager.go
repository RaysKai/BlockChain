package manager

import "github.com/linkchain/meta/block"

type BlockManager interface {
	GetBlockDetail() block.IBlock
}
