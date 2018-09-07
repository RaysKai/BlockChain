package cmd


import (
	"crypto/sha256"

	"github.com/spf13/cobra"
	"github.com/linkchain/poa/poamanager"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/poa/meta"
	"github.com/linkchain/common/math"
)

func init() {
	RootCmd.AddCommand(txCmd)
	txCmd.AddCommand(createTxCmd, signTxCmd, sendTxCmd)
}

var txCmd = &cobra.Command{
	Use:   "tx",
	Short: "all tx related command",
}

var createTxCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new tx",
	Run: func(cmd *cobra.Command, args []string) {
		println("New tx generated")
		fromAddress := math.Hash(sha256.Sum256([]byte("lf")))
		toAddress := math.Hash(sha256.Sum256([]byte("lc")))
		formAccount := &meta.POAAccount{AccountID:meta.POAAccountID{ID:fromAddress}}
		toAccount := &meta.POAAccount{AccountID:meta.POAAccountID{ID:toAddress}}
		amount := &meta.POAAmount{Value:10}
		tx := poamanager.GetManager().TransactionManager.NewTransaction(formAccount,toAccount,amount)
		log.Info("createtx","data",tx)
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
