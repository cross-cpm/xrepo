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
	case "up":
		doUpCmd(extfile)
	default:
		dumpUsage()
	}
}

func doUpCmd(extfile string) {
	log.Println("up")

	externals := repo.NewExternals(extfile)
	externals.Load()
	externals.Foreach(func(url string, workdir string, info *repo.Info) {
		log.Println("update", url, "...")
		e := repo.NewExecutor(url, workdir, info)
		err := e.Update()
		if err != nil {
			log.Fatal(err)
		}
	})
}
