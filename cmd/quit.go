package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

func init() {
	RootCmd.AddCommand(quitCmd)
}

var quitCmd = &cobra.Command{
	Use:   "quit",
	Short: "quit linkchain cmd",
	Run: func(cmd *cobra.Command, args []string) {
		os.Exit(0)
	},
}
