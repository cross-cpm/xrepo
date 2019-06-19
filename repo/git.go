package repo

import (
	"os"
)

type gitExecutor struct {
	url     string
	workdir string
	info    *Info
}

func newGitExecutor(url string, workdir string, info *Info) *gitExecutor {
	return &gitExecutor{url, workdir, info}
}

func (e *gitExecutor) Update() error {
	//log.Println("workdir", e.workdir)
	//log.Println("info", e.info)
	if _, err := os.Stat(e.workdir); os.IsNotExist(err) {
		err = run_shell("git", "clone", e.url, e.workdir)
		if err != nil {
			return err
		}
	}

	err := run_shell_in_dir(e.workdir, "git", "checkout", e.info.Branch)
	if err != nil {
		return err
	}

	if e.info.Ref != "HEAD" {
		err = run_shell_in_dir(e.workdir, "git", "checkout", e.info.Ref)
		if err != nil {
			return err
		}
	}

	return nil
}
