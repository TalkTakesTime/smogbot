package main

import (
	"flag"
	"fmt"
	"github.com/TalkTakesTime/smogbot"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) > 0 {
		smogbot.Start(args[0])
	} else {
		fmt.Println("no URL given -- please give the base URL of the thread you want to get replays from")
	}
}
