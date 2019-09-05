package main

import (
	"fmt"
	"log"
)

func cliRevDiff(extfile string) {
	externals := NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *RepoInfo) {
		idx = idx + 1
		e := NewCvs(url, info)
		ref, err := e.Revision()
		if err != nil {
			log.Fatal(err)
		}
		if ref != info.Ref {
			fmt.Printf("=== [%d/%d] %s reversion:\n", idx, count, url)
			fmt.Printf("    repo ref: %s\n", ref)
			fmt.Printf("   configure: %s\n", info.Ref)
		}
	})
}

func cliRevList(extfile string) {
	externals := NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *RepoInfo) {
		idx = idx + 1
		fmt.Printf("=== [%d/%d] %s reversion:\n", idx, count, url)
		e := NewCvs(url, info)
		ref, err := e.Revision()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("    repo ref: %s\n", ref)
		fmt.Printf("   configure: %s\n", info.Ref)
	})
}

func cliRevSave(extfile string) {
	externals := NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *RepoInfo) {
		idx = idx + 1
		fmt.Printf("=== [%d/%d] check %s revision ...\n", idx, count, url)
		e := NewCvs(url, info)
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
