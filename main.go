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
		doCheckout(extfile)
	case "co":
		doCheckout(extfile)
	case "pull":
		doPull(extfile)
	case "rev":
		switch subcmd {
		case "list":
			doRevList(extfile)
		}
	default:
		dumpUsage()
	}
}

func doCheckout(extfile string) {
	externals := repo.NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *repo.Info) {
		idx = idx + 1
		log.Printf("[%d/%d] checkout %s ...\n", idx, count, url)
		e := repo.NewExecutor(url, info)
		err := e.Checkout()
		if err != nil {
			log.Fatal(err)
		}
	})
}

func doPull(extfile string) {
	externals := repo.NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *repo.Info) {
		idx = idx + 1
		log.Printf("[%d/%d] pull %s ...\n", idx, count, url)
		e := repo.NewExecutor(url, info)
		err := e.Pull()
		if err != nil {
			log.Fatal(err)
		}
	})
}

func doRevList(extfile string) {
	externals := repo.NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *repo.Info) {
		idx = idx + 1
		log.Printf("[%d/%d] %s reversion:\n", idx, count, url)
		e := repo.NewExecutor(url, info)
		ref, err := e.Revision()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("    external ref: %s\n", info.Ref)
		log.Printf("   repo real ref: %s\n", ref)
	})
}
