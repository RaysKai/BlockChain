package main

import (
	"github.com/linkchain/node"
	"os"
	"github.com/linkchain/common/util/log"
)


func main() {

	//init log
	log.Root().SetHandler(
		log.LvlFilterHandler(log.Lvl(log.LvlInfo),
			log.StreamHandler(os.Stdout, log.TerminalFormat(true))))

	node.Init();

	node.Run();
}
