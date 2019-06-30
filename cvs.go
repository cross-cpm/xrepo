package main

import "log"

type Cvs interface {
	Checkout() error
	Pull() (string, error)
	Revision() (string, error)
	Push() error
}

func NewCvs(url string, info *RepoInfo) Cvs {
	var workdir string
	for src, dsts := range info.Targets {
		if src == "./" {
			workdir = dsts[0]
		}
	}

	if info.Cvs == "" || info.Cvs == "git" {
		return NewGitCvs(workdir, url, info)
	} else {
		log.Fatal("unsupport cvs type:", info.Cvs)
	}

	return nil
}
