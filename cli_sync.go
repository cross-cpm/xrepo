package main

import (
	"fmt"
	"log"
)

func cliCheckout(extfile string) {
	externals := NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *RepoInfo) {
		idx = idx + 1
		fmt.Printf("=== [%d/%d] checkout %s ...\n", idx, count, url)
		e := NewCvs(url, info)
		err := e.Checkout()
		if err != nil {
			log.Fatal(err)
		}
	})
}

func cliPull(extfile string) {
	externals := NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *RepoInfo) {
		idx = idx + 1
		fmt.Printf("=== [%d/%d] pull %s ...\n", idx, count, url)
		e := NewCvs(url, info)
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

func cliPush(extfile string) {
	externals := NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *RepoInfo) {
		idx = idx + 1
		fmt.Printf("=== [%d/%d] push %s ...\n", idx, count, url)
		e := NewCvs(url, info)
		err := e.Push()
		if err != nil {
			log.Println(err)
		}
	})
}
