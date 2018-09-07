package cmd


import (
	"github.com/spf13/cobra"
	"github.com/linkchain/poa/poamanager"
	"github.com/linkchain/common/util/log"
	"strconv"
)

func init() {
	RootCmd.AddCommand(mineCmd)
	RootCmd.AddCommand(chainInfoCmd)
	RootCmd.AddCommand(blockCmd)
	blockCmd.AddCommand(heightCmd)
}

var mineCmd = &cobra.Command{
	Use:   "mine",
	Short: "generate a new block",
	Run: func(cmd *cobra.Command, args []string) {
		block := poamanager.GetManager().BlockManager.NewBlock()
		poamanager.GetManager().BlockManager.ProcessBlock(block)
	},
}

var chainInfoCmd = &cobra.Command{
	Use:   "chaininfo",
	Short: "getBlockChainInfo",
	Run: func(cmd *cobra.Command, args []string) {
		poamanager.GetManager().ChainManager.GetBlockChainInfo()
	},
}

var loadChainCmd = &cobra.Command{
	Use:   "loadchain",
	Short: "loadchain",
	Run: func(cmd *cobra.Command, args []string) {
		poamanager.GetManager().ChainManager.UpdateChain()
	},
}

var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "block command",
}

var heightCmd = &cobra.Command{
	Use:   "height",
	Short: "get a block by height",
	Run: func(cmd *cobra.Command, args []string) {
		example := []string{"example","block height 0"}
		if len(args) != 1 {
			log.Error("getblockbyheight","error","please input height",example[0],example[1])
			return
		}

		height,error := strconv.Atoi(args[0])
		if error != nil {
			log.Error("getblockbyheight ","error",error,example[0],example[1])
			return
		}

		if uint32(height) > poamanager.GetManager().ChainManager.GetBestBlock().GetHeight() || height < 0{
			log.Error("getblockbyheight ","error","height is out of range",example[0],example[1])
			return
		}
		log.Info("block","data",poamanager.GetManager().ChainManager.GetBlockByHeight(uint32(height)))
	},
}





