package main

import (
	"os"
	"strings"
)

type gitExecutor struct {
	url     string
	workdir string
	info    *Info
}

func newGitExecutor(workdir string, url string, info *Info) *gitExecutor {
	// TODO: check url branck ref invald
	return &gitExecutor{url, workdir, info}
}

func (e *gitExecutor) Checkout() error {
	//log.Println("workdir", e.workdir)
	//log.Println("info", e.info)
	if _, err := os.Stat(e.workdir); os.IsNotExist(err) {
		return e.initial_checkout()
	} else {
		return e.update_checkout()
	}
}

func (e *gitExecutor) initial_checkout() error {
	err := shell_run("git", "clone", e.url, e.workdir)
	if err != nil {
		return err
	}

	if e.info.Ref == "HEAD" {
		err := shell_run_in_dir(e.workdir, "git", "checkout", e.info.Branch)
		if err != nil {
			return err
		}
	} else {
		err := shell_run_in_dir(e.workdir, "git", "checkout", e.info.Ref)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *gitExecutor) update_checkout() error {
	if e.info.Ref == "HEAD" {
		err := shell_run_in_dir(e.workdir, "git", "checkout", e.info.Branch)
		if err != nil {
			return err
		}

		err = shell_run_in_dir(e.workdir, "git", "pull")
		if err != nil {
			return err
		}
	} else {
		err := shell_run_in_dir(e.workdir, "git", "fetch")
		if err != nil {
			return err
		}

		err = shell_run_in_dir(e.workdir, "git", "checkout", e.info.Ref)
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *gitExecutor) Pull() error {
	err := shell_run_in_dir(e.workdir, "git", "checkout", e.info.Branch)
	if err != nil {
		return err
	}

	err = shell_run_in_dir(e.workdir, "git", "pull")
	if err != nil {
		return err
	}

	return nil
}

func (e *gitExecutor) Revision() (string, error) {
	rev, err := shell_get_in_dir(e.workdir, "git", "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}

	return strings.Trim(rev, " \r\n"), err
}
