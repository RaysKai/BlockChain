package main

import (
	"os"
	"fmt"
	"github.com/linkchain/common/util/log"
	"github.com/linkchain/cmd"
	"bufio"
	"strings"
)


func main() {
	//init log
	log.Root().SetHandler(
		log.LvlFilterHandler(log.Lvl(log.LvlInfo),
			log.StreamHandler(os.Stdout, log.TerminalFormat(true))))


	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(">: ")
		// Scans a line from Stdin(Console)
		scanner.Scan()
		// Holds the string that scanned
		text := scanner.Text()
		if len(text) != 0 {
			words := strings.Fields(text)

			cmd.RootCmd.SetArgs(words)
			cmd.RootCmd.Execute()

		}
	}
}
