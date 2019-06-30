package main

import "log"

func cliRevList(extfile string) {
	externals := NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *RepoInfo) {
		idx = idx + 1
		log.Printf("=== [%d/%d] %s reversion:\n", idx, count, url)
		e := NewCvs(url, info)
		ref, err := e.Revision()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("    repo ref: %s\n", ref)
		log.Printf("   configure: %s\n", info.Ref)
	})
}

func cliRevSave(extfile string) {
	externals := NewExternals(extfile)
	externals.Load()
	idx := 0
	count := externals.Count()
	externals.Foreach(func(url string, info *RepoInfo) {
		idx = idx + 1
		log.Printf("=== [%d/%d] check %s revision ...\n", idx, count, url)
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
