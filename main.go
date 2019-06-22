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
		doCheckout(extfile)
	case "co":
		doCheckout(extfile)
	case "pull":
		doPull(extfile)
	case "push":
		doPush(extfile)
	case "rev":
		switch subcmd {
		case "list":
			doRevList(extfile)
		case "save":
			doRevSave(extfile)
		}
	default:
		dumpUsage()
	}
}

func doCheckout(extfile string) {
	externals := NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *Info) {
		idx = idx + 1
		log.Printf("=== [%d/%d] checkout %s ...\n", idx, count, url)
		e := NewExecutor(url, info)
		err := e.Checkout()
		if err != nil {
			log.Fatal(err)
		}
	})
}

func doPull(extfile string) {
	externals := NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *Info) {
		idx = idx + 1
		log.Printf("=== [%d/%d] pull %s ...\n", idx, count, url)
		e := NewExecutor(url, info)
		ref, err := e.Pull()
		if err != nil {
			log.Println(err)
		}
		if ref != info.Ref {
			log.Printf("set configure to ref(%s)", ref)
			info.Ref = ref
		}
	})

	err := externals.Save()
	if err != nil {
		log.Fatal(err)
	}
}

func doPush(extfile string) {
	externals := NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *Info) {
		idx = idx + 1
		log.Printf("=== [%d/%d] push %s ...\n", idx, count, url)
		e := NewExecutor(url, info)
		err := e.Push()
		if err != nil {
			log.Println(err)
		}
	})
}

func doRevList(extfile string) {
	externals := NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *Info) {
		idx = idx + 1
		log.Printf("=== [%d/%d] %s reversion:\n", idx, count, url)
		e := NewExecutor(url, info)
		ref, err := e.Revision()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("    repo ref: %s\n", ref)
		log.Printf("   configure: %s\n", info.Ref)
	})
}

func doRevSave(extfile string) {
	externals := NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *Info) {
		idx = idx + 1
		log.Printf("=== [%d/%d] check %s revision ...\n", idx, count, url)
		e := NewExecutor(url, info)
		ref, err := e.Revision()
		if err != nil {
			log.Fatal(err)
		}
		if ref != info.Ref {
			log.Printf("set configure to ref(%s)", ref)
			info.Ref = ref
		}
	})

	err := externals.Save()
	if err != nil {
		log.Fatal(err)
	}
}
