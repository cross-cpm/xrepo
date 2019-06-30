package main

import (
	"os"
	"strings"
)

type gitCvs struct {
	url     string
	workdir string
	info    *RepoInfo
}

func NewGitCvs(workdir string, url string, info *RepoInfo) *gitCvs {
	// TODO: check url branck ref invald
	return &gitCvs{url, workdir, info}
}

func (e *gitCvs) Checkout() error {
	//log.Println("workdir", e.workdir)
	//log.Println("info", e.info)
	if _, err := os.Stat(e.workdir); os.IsNotExist(err) {
		return e.initial_checkout()
	} else {
		return e.update_checkout()
	}
}

func (e *gitCvs) initial_checkout() error {
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

func (e *gitCvs) update_checkout() error {
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

func (e *gitCvs) Pull() (string, error) {
	err := shell_run_in_dir(e.workdir, "git", "checkout", e.info.Branch)
	if err != nil {
		return "", err
	}

	err = shell_run_in_dir(e.workdir, "git", "pull")
	if err != nil {
		return "", err
	}

	return e.Revision()
}

func (e *gitCvs) Revision() (string, error) {
	rev, err := shell_get_in_dir(e.workdir, "git", "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}

	return strings.Trim(rev, " \r\n"), err
}

func (e *gitCvs) Push() error {
	return shell_run_in_dir(e.workdir, "git", "push")
}
