package main

import (
	"log"
	"os"
)

func dumpUsage() {
	log.Println("xrepo usage:")
}

func main() {
	log.SetFlags(0)

	//var apppath = filepath.Dir(os.Args[0])
	var (
		extfile = "externals.json"
		cmd     string
		subcmd  string
	)

	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	if len(os.Args) > 2 {
		subcmd = os.Args[2]
	}

	log.Println("externals file:", extfile)

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
		}
	default:
		dumpUsage()
	}
}
