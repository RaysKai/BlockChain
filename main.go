package main

import (
	"os"
	"fmt"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/cmd"
)


func main() {
	//init log
	log.Root().SetHandler(
		log.LvlFilterHandler(log.Lvl(log.LvlInfo),
			log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	//command line tools
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
