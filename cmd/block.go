package cmd


import (
	"github.com/spf13/cobra"
	"github.com/linkchain/poa/poamanager"
)

func init() {
	RootCmd.AddCommand(mineCmd)
	RootCmd.AddCommand(chainInfoCmd)
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

