package cmd

import (
	"github.com/spf13/cobra"
	"github.com/linkchain/node")

func init() {
	RootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start linkchain node",
	Run: func(cmd *cobra.Command, args []string) {
		node.Init();
		node.Run();
	},
}
