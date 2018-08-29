package cmd


import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(txCmd)
	txCmd.AddCommand(createTxCmd, signTxCmd, sendTxCmd)
}

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "generate a new block",
}

var createTxCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new tx",
	Run: func(cmd *cobra.Command, args []string) {
		println("New tx generated")
	},
}

var signTxCmd = &cobra.Command{
	Use:   "sign",
	Short: "sign a new tx",
	Run: func(cmd *cobra.Command, args []string) {
		println("Tx signed")
	},
}

var sendTxCmd = &cobra.Command{
	Use:   "send",
	Short: "send a new tx to network",
	Run: func(cmd *cobra.Command, args []string) {
		println("Tx send out")
	},
}
