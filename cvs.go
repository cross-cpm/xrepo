package main

import "log"

type Cvs interface {
	Checkout() error
	Pull() (string, error)
	Push() error
	Revision() (string, error)
}

func NewCvs(url string, info *RepoInfo) Cvs {
	var workdir string
	for src, dsts := range info.Targets {
		if src == "./" {
			workdir = dsts[0]
		}
	}

	if info.Cvs == "" {
		// FIXME: parse cvs from url
		info.Cvs = "git"
	}

	switch info.Cvs {
	case "git":
		return NewGitCvs(workdir, url, info)
	default:
		log.Fatal("unsupport cvs type:", info.Cvs)
	}

	return nil
}
