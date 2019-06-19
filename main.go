package main

import (
	"log"
	"os"
	"xrepo/repo"
)

func dumpUsage() {
	log.Println("xrepo usage:")
}

func main() {
	//var apppath = filepath.Dir(os.Args[0])
	var (
		extfile = "externals.json"
		cmd     string
	)

	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	log.Println("externals file:", extfile)

	switch cmd {
	case "checkout":
		doCheckout(extfile)
	case "co":
		doCheckout(extfile)
	default:
		dumpUsage()
	}
}

func doCheckout(extfile string) {
	externals := repo.NewExternals(extfile)
	externals.Load()
	externals.Foreach(func(url string, info *repo.Info) {
		log.Println("checkout", url, "...")
		e := repo.NewExecutor(url, info)
		err := e.Checkout()
		if err != nil {
			log.Fatal(err)
		}
	})
}
