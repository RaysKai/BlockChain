package manager

import "github.com/linkchain/meta"

type BlockManager interface {
	GetBlockDetail() meta.Block
}
