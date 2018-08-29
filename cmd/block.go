package cmd


import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(mineCmd)
}

var mineCmd = &cobra.Command{
	Use:   "mine",
	Short: "generate a new block",
	Run: func(cmd *cobra.Command, args []string) {
		println("New block generated, You got 100000000 God Coin")
	},
}