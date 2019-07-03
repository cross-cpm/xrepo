package main

import (
	"fmt"
	"os"
)

func dumpUsage() {
	fmt.Print(`usage: xrepo <command> <args>

commands:
   checkout (co)  checkout every repo to current revision
   pull           update every repo to newest revision
   push           push work revision to remote repo
   rev list       list work revision of every repo
   rev save       write work revision to externals file
`)
}

func main() {
	//log.SetFlags(0)

	var (
		extfile = "externals.yaml"
		cmd     string
		subcmd  string
	)

	if _, err := os.Stat(extfile); os.IsNotExist(err) {
		extfile = "externals.json"
	}

	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	if len(os.Args) > 2 {
		subcmd = os.Args[2]
	}

	switch cmd {
	case "checkout":
		cliCheckout(extfile)
	case "co":
		cliCheckout(extfile)
	case "pull":
		cliPull(extfile)
	case "push":
		cliPush(extfile)
	case "rev":
		switch subcmd {
		case "list":
			cliRevList(extfile)
		case "save":
			cliRevSave(extfile)
		default:
			dumpUsage()
		}
	default:
		dumpUsage()
	}
}
